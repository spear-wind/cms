package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/dave-malone/email"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var (
	ErrInvalidUserId   = errors.New("Invalid user id")
	ErrUserDoesntExist = errors.New("This user doesnt exist")
)

func InitRoutes(router *mux.Router, formatter *render.Render) {
	repo := repo()
	sender := email.NewSender()

	router.HandleFunc("/user/register", userRegistrationHandler(formatter, repo, sender)).Methods("POST")
	router.HandleFunc("/user/verify/{verificationCode}", emailVerificationHandler(formatter, repo)).Methods("GET")
	router.HandleFunc("/user", createUserHandler(formatter, repo, sender)).Methods("POST")
	router.HandleFunc("/user", getUserListHandler(formatter, repo)).Methods("GET")
	router.HandleFunc("/user/{id}", getUserHandler(formatter, repo)).Methods("GET")
}

func repo() repository {
	profile := os.Getenv("PROFILE")

	var repo repository

	if profile == "mysql" {
		// db, err := common.NewDbConn()
		// if err != nil {
		// 	repo = newMysqlRepository(db)
		// }
		//TODO - what backing store will we use for this service?
	} else {
		fmt.Println("Using in-memory repositories")
		repo = newInMemoryRepository()
	}

	return repo
}

func userRegistrationHandler(formatter *render.Render, repo repository, sender email.Sender) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var user User

		if err := json.Unmarshal(payload, &user); err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse create user request")
			return
		}

		if repo.exists(user) {
			formatter.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"errors": "This user already exists",
			})
			return
		}

		if result, err := user.register(); result.HasErrors() {
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

		if err := repo.add(user); err != nil {
			formatter.JSON(w, http.StatusInternalServerError, map[string]interface{}{
				"user":  user,
				"error": err.Error(),
			})

			return
		}

		w.Header().Add("Location", fmt.Sprintf("/user/%d", user.ID))
		formatter.JSON(w, http.StatusCreated, user)

		//TODO - refactor this to use some sort of ApplicationEvents framework
		messageBody := email.NewFileBasedHTMLTemplateMessageBody("./email-templates/user-registration.tpl", user)

		emailMessage := email.NewMessage(
			"no-reply@spearwind.io",
			user.Email,
			"SpearWind.io - New Account Verification",
			messageBody,
		)

		sender.Send(emailMessage)
	}
}

func emailVerificationHandler(formatter *render.Render, repo repository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		verificationCode := vars["verificationCode"]
		user := repo.findByVerificationCode(verificationCode)

		if len(verificationCode) == 0 || user == nil {
			formatter.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"errors": "Invalid Verification Code",
			})
			return
		}

		if err := user.verify(verificationCode); err != nil {
			formatter.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"errors": err.Error(),
			})
			return
		}

		if err := repo.update(*user); err != nil {
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

func createUserHandler(formatter *render.Render, repo repository, sender email.Sender) http.HandlerFunc {
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

		if err := repo.add(user); err != nil {
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

func getUserListHandler(formatter *render.Render, repo repository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		users := repo.listUsers()

		formatter.JSON(w, http.StatusOK, map[string]interface{}{
			"users": users,
			"total": len(users),
		})
	}
}

func getUserHandler(formatter *render.Render, repo repository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		userID := vars["id"]

		if user, err := repo.getUser(userID); err != nil {
			formatter.JSON(w, http.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			formatter.JSON(w, http.StatusOK, user)
		}
	}
}
