package httphandlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/wisdommatt/creativeadvtech-assessment/components/users"
)

type userApiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	User    *users.User `json:"user"`
}

// HandleCreateUserEndpoint is the http handler for create user
// endpoint.
func HandleCreateUserEndpoint(userService users.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var payload users.User
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(userApiResponse{
				Status:  "error",
				Message: "invalid json payload",
			})
			return
		}
		user, err := userService.CreateUser(r.Context(), payload)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(rw).Encode(userApiResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(userApiResponse{
			Status:  "success",
			Message: "user created successfully",
			User:    user,
		})
	}
}

// HandleGetUserEndpoint is the http request handler to get user details.
func HandleGetUserEndpoint(userService users.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userId")
		user, err := userService.GetUser(r.Context(), userID)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(rw).Encode(userApiResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(userApiResponse{
			Status:  "success",
			Message: "user retrieved successfully",
			User:    user,
		})
	}
}
