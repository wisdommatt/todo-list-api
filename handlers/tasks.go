package httphandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/wisdommatt/todo-list-api/services/tasks"
	"github.com/wisdommatt/todo-list-api/services/users"
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

// HandleCreateTaskEndpoint is the http endpoint handler for creating a new task.
func HandleCreateTaskEndpoint(tasksService *tasks.Service, usersService *users.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var payload tasks.Task
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			ErrorResponse(rw, "error", "invalid json payload", http.StatusBadRequest)
			return
		}
		_, err = usersService.GetUser(r.Context(), payload.UserID)
		if err != nil {
			ErrorResponse(rw, "error", "user does not exist", http.StatusBadRequest)
			return
		}
		// checking if the new task is overlapping with another existing task.
		overlappingTask, _ := tasksService.GetTaskWithinTimeRange(r.Context(), payload.UserID, payload.StartTime, payload.EndTime)
		if overlappingTask != nil {
			errMsg := fmt.Sprintf("this task if overlapping with %s, pick another time", overlappingTask.Title)
			ErrorResponse(rw, "error", errMsg, http.StatusBadRequest)
			return
		}
		task, err := tasksService.CreateTask(r.Context(), payload)
		if err != nil {
			ErrorResponse(rw, "error", errSomethingWentWrongMsg, http.StatusInternalServerError)
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
func HandleGetTaskEndpoint(tasksService *tasks.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		taskID := chi.URLParam(r, "taskId")
		task, err := tasksService.GetTask(r.Context(), taskID)
		if err != nil {
			ErrorResponse(rw, "error", "task does not exist", http.StatusBadRequest)
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

// HandleGetTasksEndpoint is the http endpoint handler for retrieving user tasks.
func HandleGetTasksEndpoint(tasksService *tasks.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userId")
		lastID := r.URL.Query().Get("lastId")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		tasks, err := tasksService.GetTasks(r.Context(), userID, lastID, limit)
		if err != nil {
			ErrorResponse(rw, "error", errSomethingWentWrongMsg, http.StatusInternalServerError)
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
func HandleDeleteTaskEndpoint(tasksService *tasks.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		taskID := chi.URLParam(r, "taskId")
		_, err := tasksService.GetTask(r.Context(), taskID)
		if err != nil {
			ErrorResponse(rw, "error", "invalid task id", http.StatusBadRequest)
			return
		}
		task, err := tasksService.DeleteTask(r.Context(), taskID)
		if err != nil {
			ErrorResponse(rw, "error", errSomethingWentWrongMsg, http.StatusInternalServerError)
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
func HandleUpdateTaskEndpoint(tasksService *tasks.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		taskID := chi.URLParam(r, "taskId")
		var payload updateTaskPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			ErrorResponse(rw, "error", "invalid json payload", http.StatusBadRequest)
			return
		}
		task, err := tasksService.UpdateTask(r.Context(), taskID, tasks.Task{
			Status: payload.Status,
		})
		if err != nil {
			ErrorResponse(rw, "error", errSomethingWentWrongMsg, http.StatusInternalServerError)
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
