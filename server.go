package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/dave-malone/email"
	"github.com/gorilla/mux"
	"github.com/spear-wind/cms/auth"
	"github.com/spear-wind/cms/events"
	"github.com/spear-wind/cms/user"
	"github.com/unrolled/render"
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if awsEndpoint != "" && awsAccessKeyID != "" && awsSecretAccessKey != "" {
		fmt.Println("Using Amazon SES Email Sender")
		email.NewSender = email.NewAmazonSESSender(awsEndpoint, awsAccessKeyID, awsSecretAccessKey)
	} else {
		email.NewSender = email.NewNoopSender
	}

	eventPublisher := events.NewSynchEventPublisher()
	eventPublisher.Add(events.NewEmailEventSubscriber(email.NewSender()))

	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	auth.InitRoutes(mx, formatter)
	user.InitRoutes(mx, formatter, eventPublisher)
	n.UseHandler(mx)
	return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/test", testHandler(formatter)).Methods("GET")
	mx.Handle("/securetest", wrapHandler(formatter, testSecureHandler(formatter)))
}

func testHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"This is a test"})
	}
}

func testSecureHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"This is a secured test"})
	}
}

func wrapHandler(formatter *render.Render, handler http.HandlerFunc) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(auth.IsAuthorized(formatter)),
		negroni.Wrap(handler),
	)
}
