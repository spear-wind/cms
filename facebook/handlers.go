package facebook

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spear-wind/cms/auth"
	"github.com/spear-wind/cms/user"
	"github.com/unrolled/render"
)

func InitRoutes(router *mux.Router, formatter *render.Render, userRepository user.UserRepository, fbClient *Client) {
	router.HandleFunc("/facebook/login", facebookLoginHandler(formatter, userRepository, fbClient)).Methods("POST")
}

func facebookLoginHandler(formatter *render.Render, userRepository user.UserRepository, fbClient *Client) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)

		var cmd loginCommand

		if err := json.Unmarshal(payload, &cmd); err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse fb auth response: "+err.Error())
			return
		}

		fbUser, err := fbClient.getUser(cmd)

		if err != nil {
			formatter.Text(w, http.StatusBadRequest, err.Error())
			return
		}

		//TODO look up user by fb id
		//TODO - look up user by email
		//TODO - save new or update existing user

		tokenString, err := auth.GenerateToken()
		if err != nil {
			formatter.JSON(w, http.StatusOK, struct{ Message string }{err.Error()})
			return
		}

		data := struct {
			User  *user.User
			Token string
		}{
			fbUser,
			tokenString,
		}

		formatter.JSON(w, http.StatusOK, data)
	}
}
