package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spear-wind/cms/events"
	"github.com/unrolled/render"
)

var (
	errInvalidUserID   = errors.New("Invalid user id")
	errUserDoesntExist = errors.New("This user doesnt exist")
)

func InitRoutes(router *mux.Router, formatter *render.Render, userRepository UserRepository, eventPublisher events.EventPublisher) {
	router.HandleFunc("/user", createUserHandler(formatter, userRepository, eventPublisher)).Methods("POST")
	router.HandleFunc("/user", getUserListHandler(formatter, userRepository)).Methods("GET")
	router.HandleFunc("/user/{id}", getUserHandler(formatter, userRepository)).Methods("GET")
}

func createUserHandler(formatter *render.Render, userRepository UserRepository, eventPublisher events.EventPublisher) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var user User

		err := json.Unmarshal(payload, &user)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse create user request")
			return
		}

		if result := user.validate(); result.HasErrors() {
			formatter.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"errors": result.Errors,
			})
			return
		}

		if err := userRepository.Add(&user); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, map[string]interface{}{
				"user":  user,
				"error": err.Error(),
			})
		} else {
			w.Header().Add("Location", fmt.Sprintf("/user/%d", user.ID))
			formatter.JSON(w, http.StatusCreated, user)
			//TODO newUserCreatedEvent(user)
		}
	}
}

func getUserListHandler(formatter *render.Render, userRepository UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		users := userRepository.listUsers()

		formatter.JSON(w, http.StatusOK, map[string]interface{}{
			"users": users,
			"total": len(users),
		})
	}
}

func getUserHandler(formatter *render.Render, userRepository UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		userID := vars["id"]

		if user, err := userRepository.getUser(userID); err != nil {
			formatter.JSON(w, http.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			formatter.JSON(w, http.StatusOK, user)
		}
	}
}
