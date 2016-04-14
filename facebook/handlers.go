package facebook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spear-wind/cms/auth"
	"github.com/spear-wind/cms/user"
	"github.com/unrolled/render"
)

func InitRoutes(router *mux.Router, formatter *render.Render, userRepository user.UserRepository, fbClient Client) {
	router.HandleFunc("/facebook/login", facebookLoginHandler(formatter, userRepository, fbClient)).Methods("POST")
}

func facebookLoginHandler(formatter *render.Render, userRepository user.UserRepository, fbClient Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)

		var cmd loginCommand

		if err := json.Unmarshal(payload, &cmd); err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse fb auth response: "+err.Error())
			return
		}

		if result := cmd.validate(); result.HasErrors() {
			formatter.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"errors": result.Errors,
			})
			return
		}

		fbUser, err := fbClient.getUser(cmd)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, err.Error())
			return
		}

		fmt.Printf("fb user: %v\n", fbUser)

		existingUser := userRepository.FindByFacebookID(fbUser.FacebookID)
		if existingUser == nil {
			existingUser = userRepository.FindByEmail(fbUser.Email)

			if existingUser != nil {
				existingUser.FacebookID = fbUser.FacebookID
				if err := userRepository.Update(existingUser); err != nil {
					formatter.JSON(w, http.StatusOK, struct{ Message string }{err.Error()})
					return
				}
			} else {
				if err := userRepository.Add(fbUser); err != nil {
					formatter.JSON(w, http.StatusOK, struct{ Message string }{err.Error()})
					return
				}

				existingUser = fbUser
			}
		}

		tokenString, err := auth.GenerateToken(existingUser.ID)
		if err != nil {
			formatter.JSON(w, http.StatusOK, struct{ Message string }{err.Error()})
			return
		}

		data := struct {
			User  *user.User
			Token string
		}{
			existingUser,
			tokenString,
		}

		formatter.JSON(w, http.StatusOK, data)
	}
}
