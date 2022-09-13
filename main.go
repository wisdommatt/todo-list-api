package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	handlers "github.com/wisdommatt/todo-list-api/handlers"
	"github.com/wisdommatt/todo-list-api/internal/jwt"
	"github.com/wisdommatt/todo-list-api/services/tasks"
	"github.com/wisdommatt/todo-list-api/services/users"
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
	usersService := users.NewUsersService(mongoDB, log)
	tasksService := tasks.NewService(usersService, mongoDB, log)

	router := chi.NewRouter()
	router.Route("/users/", func(r chi.Router) {
		r.Post("/", handlers.HandleCreateUserEndpoint(usersService))
		r.Post("/login", handlers.HandleUserLoginEndpoint(usersService))

		r.Group(func(r chi.Router) {
			r.Use(isLoggedInMiddleware)
			r.Get("/{userId}", handlers.HandleGetUserEndpoint(usersService))
			r.Get("/", handlers.HandleGetUsersEndpoint(usersService))
			r.Delete("/{userId}", handlers.HandleDeleteUserEndpoint(usersService))
			r.Get("/{userId}/tasks", handlers.HandleGetTasksEndpoint(tasksService))
		})
	})

	router.Route("/tasks/", func(r chi.Router) {
		r.Use(isLoggedInMiddleware)
		r.Post("/", handlers.HandleCreateTaskEndpoint(tasksService, usersService))
		r.Get("/{taskId}", handlers.HandleGetTaskEndpoint(tasksService))
		r.Put("/{taskId}", handlers.HandleUpdateTaskEndpoint(tasksService))
		r.Delete("/{taskId}", handlers.HandleDeleteTaskEndpoint(tasksService))
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

func isLoggedInMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get("Authorization")
		authToken = strings.ReplaceAll(authToken, "Bearer ", "")
		_, err := jwt.Decode([]byte(os.Getenv("JWT_SECRET")), authToken)
		if err != nil {
			json.NewEncoder(rw).Encode(map[string]string{
				"status":  "unauthorized",
				"message": "you are not authorized to proceed",
			})
			return
		}
		h.ServeHTTP(rw, r)
	})
}
