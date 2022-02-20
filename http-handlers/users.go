package httphandlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/wisdommatt/creativeadvtech-assessment/components/users"
)

type userApiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	User    *users.User `json:"user"`
}

type getUsersResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Users   []users.User `json:"users"`
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

// HandleGetUserEndpoint is the http endpoint handler to get user details.
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

// HandleGetUsersEndpoint is the http endpoint handler for retrieving
// users.
func HandleGetUsersEndpoint(userService users.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		lastID := r.URL.Query().Get("lastId")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		users, err := userService.GetUsers(r.Context(), lastID, limit)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(getUsersResponse{
				Status:  "error",
				Message: err.Error(),
				Users:   users,
			})
			return
		}
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(getUsersResponse{
			Status:  "success",
			Message: "users retrieved successfully",
			Users:   users,
		})
	}
}

// HandleDeleteUserEndpoint is the http endpoint handler for deleting user.
func HandleDeleteUserEndpoint(userService users.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userId")
		user, err := userService.DeleteUser(r.Context(), userID)
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
			Message: "user deleted successfully",
			User:    user,
		})
	}
}
