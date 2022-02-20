package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/wisdommatt/creativeadvtech-assessment/components/users"
	httphandlers "github.com/wisdommatt/creativeadvtech-assessment/http-handlers"
)

var defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)

	userService := users.NewService(nil, log)

	router := chi.NewRouter()
	router.Route("/users/", func(r chi.Router) {
		r.Post("/", httphandlers.HandleCreateUserEndpoint(userService))
		r.Get("/users/{userId}", httphandlers.HandleGetUserEndpoint(userService))
	})

	server := &http.Server{
		Addr:         ":" + defaultPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Infof("app running on port: %s", port)
	log.Fatal(server.ListenAndServe())
}
