package tasks

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/wisdommatt/creativeadvtech-assessment/components/users"
)

// Service is the interface that describes a task repository
// object.
type Service interface {
	CreateTask(ctx context.Context, task Task) (*Task, error)
	GetTask(ctx context.Context, taskID string) (*Task, error)
	GetTasks(ctx context.Context, userID, lastID string, limit int) ([]Task, error)
	DeleteTask(ctx context.Context, taskID string) (*Task, error)
}

// taskService is the default implementation for Service interface.
type taskService struct {
	userService users.Service
	taskRepo    repository
	log         *logrus.Logger
}

var (
	errSomethingWentWrong = fmt.Errorf("an error occured, please try again later")
)

// NewService creates a new task service.
func NewService(userService users.Service, taskRepo repository, log *logrus.Logger) *taskService {
	return &taskService{
		userService: userService,
		taskRepo:    taskRepo,
		log:         log,
	}
}

func (s *taskService) CreateTask(ctx context.Context, task Task) (*Task, error) {
	log := s.log.WithContext(ctx).WithField("task", task)
	// checking is task user is valid.
	_, err := s.userService.GetUser(ctx, task.UserID)
	if err != nil {
		return nil, err
	}
	newTask, err := s.taskRepo.saveTask(ctx, task)
	if err != nil {
		log.WithError(err).Error("an error occured while creating task")
		return nil, errSomethingWentWrong
	}
	return newTask, nil
}

func (s *taskService) GetTask(ctx context.Context, taskID string) (*Task, error) {
	log := s.log.WithContext(ctx).WithField("taskId", taskID)
	task, err := s.taskRepo.getTaskByID(ctx, taskID)
	if err != nil {
		log.WithError(err).Error("an error occured while retrieving task")
		return nil, fmt.Errorf("task does not exist")
	}
	return task, nil
}

func (s *taskService) GetTasks(ctx context.Context, userID, lastID string, limit int) ([]Task, error) {
	log := s.log.WithContext(ctx).WithField("userId", userID).WithField("lastId", lastID).
		WithField("limit", limit)
	_, err := s.userService.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	tasks, err := s.taskRepo.getTasks(ctx, userID, lastID, limit)
	if err != nil {
		log.WithError(err).Error("an error occured while retrieving tasks")
		return nil, errSomethingWentWrong
	}
	return tasks, nil
}

func (s *taskService) DeleteTask(ctx context.Context, taskID string) (*Task, error) {
	log := s.log.WithContext(ctx).WithField("taskId", taskID)
	// checking if task id is valid.
	_, err := s.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}
	task, err := s.taskRepo.deleteTaskByID(ctx, taskID)
	if err != nil {
		log.WithError(err).Error("an error occured while deleting task")
		return nil, errSomethingWentWrong
	}
	return task, nil
}
