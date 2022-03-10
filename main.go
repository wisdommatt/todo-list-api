package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/wisdommatt/creativeadvtech-assessment/components/tasks"
	"github.com/wisdommatt/creativeadvtech-assessment/components/users"
	httphandlers "github.com/wisdommatt/creativeadvtech-assessment/http-handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	godotenv.Load(".env", ".env-defaults")

	mongoDB := mustConnectMongoDB(log)
	userService := users.NewService(mongoDB, log)
	taskServie := tasks.NewService(userService, mongoDB, log)

	router := chi.NewRouter()
	router.Route("/users/", func(r chi.Router) {
		r.Post("/", httphandlers.HandleCreateUserEndpoint(userService))
		r.Get("/{userId}", httphandlers.HandleGetUserEndpoint(userService))
		r.Get("/", httphandlers.HandleGetUsersEndpoint(userService))
		r.Delete("/{userId}", httphandlers.HandleDeleteUserEndpoint(userService))
		r.Get("/{userId}/tasks", httphandlers.HandleGetTasksEndpoint(taskServie))
		r.Post("/login", httphandlers.HandleUserLoginEndpoint(userService))
	})

	router.Route("/tasks/", func(r chi.Router) {
		r.Post("/", httphandlers.HandleCreateTaskEndpoint(taskServie))
		r.Get("/{taskId}", httphandlers.HandleGetTaskEndpoint(taskServie))
		r.Delete("/{taskId}", httphandlers.HandleDeleteTaskEndpoint(taskServie))
	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Infof("app running on port: %s", port)
	log.Fatal(server.ListenAndServe())
}

func mustConnectMongoDB(log *logrus.Logger) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.WithError(err).Fatal("Unable to connect to mongodb")
	}
	return client.Database(os.Getenv("MONGODB_DATABASE_NAME"))
}
