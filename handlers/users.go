package httphandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/wisdommatt/todo-list-api/services/users"
)

type createUserInput struct {
	FirstName string `json:"firstName" bson:"firstName,omitempty"`
	LastName  string `json:"lastName" bson:"lastName,omitempty"`
	Email     string `json:"email" bson:"email,omitempty"`
	Password  string `json:"password" bson:"password,omitempty"`
}

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

type loginUserInput struct {
	Email    string `json:"email" bson:"email,omitempty"`
	Password string `json:"password" bson:"password,omitempty"`
}

type loginUserResponse struct {
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	User      *users.User `json:"user"`
	AuthToken string      `json:"authToken"`
}

func HandleCreateUserEndpoint(usersService *users.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var payload createUserInput
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			ErrorResponse(rw, "error", "invalid json payload", http.StatusBadRequest)
			return
		}
		userWithEmail, _ := usersService.GetUserByEmail(r.Context(), payload.Email)
		if userWithEmail != nil {
			errMsg := fmt.Sprintf("user with email %s already exist", payload.Email)
			ErrorResponse(rw, "error", errMsg, http.StatusBadRequest)
			return
		}
		user, err := usersService.CreateUser(r.Context(), users.User{
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Email:     payload.Email,
			Password:  payload.Password,
		})
		if err != nil {
			ErrorResponse(rw, "error", errSomethingWentWrongMsg, http.StatusInternalServerError)
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
func HandleGetUserEndpoint(usersService *users.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userId")
		user, err := usersService.GetUser(r.Context(), userID)
		if err != nil {
			ErrorResponse(rw, "error", "user does not exist", http.StatusBadRequest)
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

// HandleGetUsersEndpoint is the http endpoint handler for retrieving users.
func HandleGetUsersEndpoint(usersService *users.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		lastID := r.URL.Query().Get("lastId")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		users, err := usersService.GetUsers(r.Context(), lastID, limit)
		if err != nil {
			ErrorResponse(rw, "error", errSomethingWentWrongMsg, http.StatusInternalServerError)
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
func HandleDeleteUserEndpoint(usersService *users.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userId")
		_, err := usersService.GetUser(r.Context(), userID)
		if err != nil {
			ErrorResponse(rw, "error", "user does not exist", http.StatusBadRequest)
			return
		}
		user, err := usersService.DeleteUser(r.Context(), userID)
		if err != nil {
			ErrorResponse(rw, "error", errSomethingWentWrongMsg, http.StatusInternalServerError)
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

// HandleUserLoginEndpoint is the http endpoint handler for user login.
func HandleUserLoginEndpoint(usersService *users.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var payload loginUserInput
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			ErrorResponse(rw, "error", "invalid json payload", http.StatusBadRequest)
			return
		}
		user, authToken, err := usersService.LoginUser(r.Context(), payload.Email, payload.Password)
		if err != nil {
			ErrorResponse(rw, "error", errSomethingWentWrongMsg, http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(loginUserResponse{
			Status:    "success",
			Message:   "user login successfully",
			User:      user,
			AuthToken: authToken,
		})
	}
}
