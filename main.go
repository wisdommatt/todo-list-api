package main

import (
	"net/http"
	"os"
	"time"

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

	userService := users.NewService(nil, log)

	mux := http.NewServeMux()
	mux.Handle("/users/", httphandlers.HandleCreateUserEndpoint(userService))

	server := &http.Server{
		Addr:         ":" + defaultPort,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Infof("app running on port: %s", port)
	log.Fatal(server.ListenAndServe())
}
