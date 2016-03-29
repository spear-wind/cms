package auth

import (
	"fmt"
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
	router.HandleFunc("/verify", createVerifyHandler(formatter)).Methods("POST")
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

func createVerifyHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {

			fmt.Printf("token --- \nvalid: %v\nmethod: %v\nclaims: %v", token.Valid, token.Method, token.Claims)

			// Don't forget to validate the alg is what you expect:
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return signingKey, nil
		})

		if err == nil && token.Valid {
			formatter.JSON(w, http.StatusOK, struct{ Message string }{"Valid Token"})
		} else {
			formatter.JSON(w, http.StatusOK, struct{ Message string }{"Bad Token: " + err.Error()})
		}
	}
}
