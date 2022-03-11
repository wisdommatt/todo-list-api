package httphandlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/wisdommatt/todo-list-api/components/tasks"
)

type taskApiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Task    *tasks.Task `json:"task"`
}

type getTasksResponse struct {
	Status  string       `json:"status"`
	Message string       `json:"message"`
	Tasks   []tasks.Task `json:"tasks"`
}

type updateTaskPayload struct {
	Status string `json:"status"`
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

// HandleGetTaskEndpoint is the http endpoint handler to get task details.
func HandleGetTaskEndpoint(taskService tasks.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		taskID := chi.URLParam(r, "taskId")
		task, err := taskService.GetTask(r.Context(), taskID)
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
			Message: "task retrieved successfully",
			Task:    task,
		})
	}
}

// HandleGetTasksEndpoint is the http endpoint handler for retrieving
// user tasks.
func HandleGetTasksEndpoint(taskService tasks.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userId")
		lastID := r.URL.Query().Get("lastId")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		tasks, err := taskService.GetTasks(r.Context(), userID, lastID, limit)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(getTasksResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(getTasksResponse{
			Status:  "success",
			Message: "tasks retrieved successfully",
			Tasks:   tasks,
		})
	}
}

// HandleDeleteTaskEndpoint is the http endpoint handler for deleting task.
func HandleDeleteTaskEndpoint(taskService tasks.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		taskID := chi.URLParam(r, "taskId")
		task, err := taskService.DeleteTask(r.Context(), taskID)
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
			Message: "task deleted successfully",
			Task:    task,
		})
	}
}

// HandleUpdateTaskEndpoint is the http endpoint handler for task update.
func HandleUpdateTaskEndpoint(taskService tasks.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		taskID := chi.URLParam(r, "taskId")
		var payload updateTaskPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(rw).Encode(taskApiResponse{
				Status:  "error",
				Message: "invalid json payload",
			})
			return
		}
		task, err := taskService.UpdateTask(r.Context(), taskID, tasks.Task{
			Status: payload.Status,
		})
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
			Message: "task updated successfully",
			Task:    task,
		})
	}
}
