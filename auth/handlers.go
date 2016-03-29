package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var (
	signingKey []byte = []byte("Rsw!MPC60$dCF$*jK%0R")
)

func InitRoutes(router *mux.Router, formatter *render.Render) {
	router.HandleFunc("/login", createLoginHandler(formatter)).Methods("POST")
}

func createLoginHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Create the token
		token := jwt.New(jwt.SigningMethodHS256)
		// Set some claims
		token.Claims["foo"] = "bar"
		token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		// Sign and get the complete encoded token as a string
		tokenString, err := token.SignedString(signingKey)

		if err != nil {
			formatter.JSON(w, http.StatusOK, struct{ Message string }{err.Error()})
			return
		}

		formatter.JSON(w, http.StatusOK, struct{ Token string }{tokenString})
	}
}
