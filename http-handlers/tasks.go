package httphandlers

import (
	"encoding/json"
	"net/http"

	"github.com/wisdommatt/creativeadvtech-assessment/components/tasks"
)

type taskApiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Task    *tasks.Task `json:"task"`
}

// HandleCreateTaskEndpoint is the http endpoint handler for creating a
// new task.
func HandleCreateTaskEndpoint(taskService tasks.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var payload tasks.Task
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(taskApiResponse{
				Status:  "error",
				Message: "invalid json payload",
			})
			return
		}
		task, err := taskService.CreateTask(r.Context(), payload)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(rw).Encode(taskApiResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(taskApiResponse{
			Status:  "success",
			Message: "task created successfully",
			Task:    task,
		})
	}
}
