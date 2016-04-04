package registration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spear-wind/cms/events"
	"github.com/spear-wind/cms/user"
	"github.com/unrolled/render"
)

func InitRoutes(router *mux.Router, formatter *render.Render, userRepository user.UserRepository, eventPublisher events.EventPublisher) {
	router.HandleFunc("/register", userRegistrationHandler(formatter, userRepository, eventPublisher)).Methods("POST")
	router.HandleFunc("/verify/{verificationCode}", userVerificationHandler(formatter, userRepository)).Methods("POST")
}

func userRegistrationHandler(formatter *render.Render, userRepository user.UserRepository, eventPublisher events.EventPublisher) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var user user.User

		if err := json.Unmarshal(payload, &user); err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse create user request")
			return
		}

		if userRepository.Exists(user) {
			formatter.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"errors": "This user already exists",
			})
			return
		}

		if result, err := user.Register(); result.HasErrors() {
			formatter.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"errors": result.Errors,
			})
			return
		} else if err != nil {
			formatter.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"errors": err.Error(),
			})
			return
		}

		if err := userRepository.Add(user); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, map[string]interface{}{
				"user":  user,
				"error": err.Error(),
			})

			return
		}

		w.Header().Add("Location", fmt.Sprintf("/user/%d", user.ID))
		formatter.JSON(w, http.StatusCreated, user)

		eventPublisher.Publish(events.NewUserRegistrationEvent(user.Email, user.VerificationCode))
		fmt.Printf("New user registration event published; verification code: %s", user.VerificationCode)
	}
}

func userVerificationHandler(formatter *render.Render, userRepository user.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		verificationCode := vars["verificationCode"]
		user := userRepository.FindByVerificationCode(verificationCode)

		if len(verificationCode) == 0 || user == nil {
			formatter.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"errors": "Invalid Verification Code",
			})
			return
		}

		if err := user.Verify(verificationCode); err != nil {
			formatter.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"errors": err.Error(),
			})
			return
		}

		if err := userRepository.Update(*user); err != nil {
			formatter.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"errors": err.Error(),
			})
			return
		}

		formatter.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"success": "Your account is now verified",
		})
	}
}
