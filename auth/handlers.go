package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/spear-wind/cms/user"
	"github.com/unrolled/render"
)

var (
	signingKey = []byte("Rsw!MPC60$dCF$*jK%0R")
)

func InitRoutes(router *mux.Router, formatter *render.Render, userRepository user.UserRepository) {
	router.HandleFunc("/login", loginHandler(formatter, userRepository)).Methods("POST")
}

func IsAuthorized(formatter *render.Render) negroni.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		token, err := parseToken(req)

		if err != nil || token == nil {
			formatter.JSON(w, http.StatusUnauthorized, struct{ Error string }{"Unauthorized."})
		} else if err == nil && token.Valid != true {
			formatter.JSON(w, http.StatusUnauthorized, struct{ Message string }{"Invalid Token."})
		} else {
			next(w, req)
		}
	}
}

func loginHandler(formatter *render.Render, userRepository user.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		email := req.FormValue("email")
		password := req.FormValue("password")

		if email == "" || password == "" {
			formatter.Text(w, http.StatusBadRequest, "Username and Password are required")
			return
		}

		user := userRepository.FindByEmail(email)
		if user == nil {
			formatter.Text(w, http.StatusNotFound, "User Not Found")
			return
		}

		success, newHash := user.Authenticate(password)
		if success != true {
			formatter.Text(w, http.StatusUnauthorized, "Unauthorized.")
			return
		}

		if newHash {
			fmt.Println("Call to user.Authenticate resulted in newHash == true; we need to update this in the DB or next auth attempt will fail")
		}

		tokenString, err := GenerateToken(user.ID)
		if err != nil {
			formatter.JSON(w, http.StatusOK, struct{ Message string }{err.Error()})
			return
		}

		formatter.JSON(w, http.StatusOK, struct{ Token string }{tokenString})
	}
}

func GenerateToken(userID int64) (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims["iat"] = time.Now().Unix()
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	token.Claims["sub"] = userID
	token.Claims["iss"] = "https://cms.spearwind.io"
	token.Claims["aud"] = "TODO"
	// token.Claims["nbf"] = ""
	// token.Claims["jti"] = ""

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(signingKey)

	return tokenString, err
}

func parseToken(req *http.Request) (*jwt.Token, error) {
	token, err := jwt.ParseFromRequest(req, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return signingKey, nil
	})

	return token, err
}
