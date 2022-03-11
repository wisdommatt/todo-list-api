package httphandlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/wisdommatt/todo-list-api/components/users"
	mockUsers "github.com/wisdommatt/todo-list-api/mocks/components/users"
)

// tests for other endpoints in this application would follow this
// same pattern.
func TestHandleCreateUserEndpoint(t *testing.T) {
	userService := &mockUsers.Service{}
	userService.On("CreateUser", mock.Anything, users.User{FirstName: "Hello", LastName: "World"}).
		Return(nil, fmt.Errorf("invalid user"))
	userService.On("CreateUser", mock.Anything, users.User{FirstName: "Valid", LastName: "User"}).
		Return(&users.User{FirstName: "Valid", LastName: "User"}, nil)

	tests := []struct {
		name           string
		payload        interface{}
		wantStatusCode int
	}{
		{
			name:           "invalid payload",
			payload:        "invalid",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "CreateUser user service implementation with error",
			payload:        users.User{FirstName: "Hello", LastName: "World"},
			wantStatusCode: http.StatusInternalServerError,
		},
		{
			name:           "CreateUser user service implementation without errors",
			payload:        users.User{FirstName: "Valid", LastName: "User"},
			wantStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyJSON, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer(bodyJSON))
			rr := httptest.NewRecorder()
			handler := HandleCreateUserEndpoint(userService)
			handler.ServeHTTP(rr, req)
			if rr.Code != tt.wantStatusCode {
				t.Errorf("want %v got %v", tt.wantStatusCode, rr.Code)
			}
			if rr.Code == http.StatusOK {
				var response userApiResponse
				err := json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Fatal(err)
				}
				if response.User == nil {
					t.Errorf("want user got %v nil", response.User)
				}
			}
		})
	}
}
