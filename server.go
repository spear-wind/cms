package main

import (
	"fmt"
	"os"

	"github.com/cloudnativego/cfmgo"
	"github.com/codegangsta/negroni"
	"github.com/dave-malone/email"
	"github.com/gorilla/mux"
	"github.com/spear-wind/cms/auth"
	"github.com/spear-wind/cms/events"
	"github.com/spear-wind/cms/facebook"
	"github.com/spear-wind/cms/registration"
	"github.com/spear-wind/cms/site"
	"github.com/spear-wind/cms/user"
	"github.com/unrolled/render"
)

// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := newFormatter()
	emailSender := newEmailSender()
	eventPublisher := newEventPublisher(emailSender)
	userRepository := newUserRepository()
	facebookClient := newFacebookClient()
	siteRepository := newSiteRepository()

	n := negroni.Classic()
	router := mux.NewRouter()

	auth.InitRoutes(router, formatter, userRepository)
	registration.InitRoutes(router, formatter, userRepository, eventPublisher)
	facebook.InitRoutes(router, formatter, userRepository, facebookClient)

	userRouter := mux.NewRouter()
	user.InitRoutes(userRouter, formatter, userRepository, eventPublisher)
	router.PathPrefix("/user").Handler(negroni.New(
		negroni.HandlerFunc(auth.IsAuthorized(formatter)),
		negroni.Wrap(userRouter),
	))

	siteRouter := mux.NewRouter()
	site.InitRoutes(siteRouter, formatter, siteRepository, eventPublisher)
	router.PathPrefix("/site").Handler(negroni.New(
		negroni.HandlerFunc(auth.IsAuthorized(formatter)),
		negroni.Wrap(siteRouter),
	))

	n.UseHandler(router)
	return n
}

func newFormatter() *render.Render {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	return formatter
}

func newEmailSender() email.Sender {
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if awsEndpoint != "" && awsAccessKeyID != "" && awsSecretAccessKey != "" {
		fmt.Println("Using Amazon SES Email Sender")
		email.NewSender = email.NewAmazonSESSender(awsEndpoint, awsAccessKeyID, awsSecretAccessKey)
	} else {
		email.NewSender = email.NewNoopSender
	}

	return email.NewSender()
}

func newEventPublisher(emailSender email.Sender) events.EventPublisher {
	eventPublisher := events.NewSynchEventPublisher()
	eventPublisher.Add(events.NewEmailEventSubscriber(emailSender))
	return eventPublisher
}

func newUserRepository() user.UserRepository {
	mongoDBURL := os.Getenv("MONGO_URL")

	var repo user.UserRepository

	if len(mongoDBURL) != 0 {
		userCollection := cfmgo.Connect(cfmgo.NewCollectionDialer, mongoDBURL, "users")
		fmt.Println("Using to MongoDB user repository")
		repo = user.NewMongoUserRepository(userCollection)
	} else {
		fmt.Println("Using in-memory user repository")
		repo = user.NewInMemoryRepository()
	}

	return repo
}

func newFacebookClient() facebook.Client {
	appID := os.Getenv("FB_APP_ID")
	appSecret := os.Getenv("FB_APP_SECRET")

	return facebook.NewClient(appID, appSecret)
}

func newSiteRepository() site.SiteRepository {
	fmt.Println("Using in-memory site repository")
	return site.NewInMemoryRepository()

	// mongoDBURL := os.Getenv("MONGO_URL")
	//
	// var repo site.SiteRepository
	//
	// if len(mongoDBURL) != 0 {
	// 	siteCollection := cfmgo.Connect(cfmgo.NewCollectionDialer, mongoDBURL, "sites")
	// 	fmt.Println("Using to MongoDB site repository")
	// 	repo = site.NewMongoSiteRepository(siteCollection)
	// } else {
	// 	fmt.Println("Using in-memory site repository")
	// 	repo = site.NewInMemoryRepository()
	// }
	//
	// return repo
}
