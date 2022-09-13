package tasks

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wisdommatt/todo-list-api/services/users"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Task struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Title     string    `json:"title" bson:"title,omitempty"`
	StartTime time.Time `json:"startTime" bson:"startTime,omitempty"`
	EndTime   time.Time `json:"endTime" bson:"endTime,omitempty"`
	UserID    string    `json:"userId" bson:"userId,omitempty"`
	Status    string    `json:"status" bson:"status,omitempty"`
	TimeAdded time.Time `json:"-" bson:"timeAdded,omitempty"`
}

type Service struct {
	usersService *users.Service
	dbCollection *mongo.Collection
	log          *logrus.Logger
}

func NewService(usersService *users.Service, db *mongo.Database, log *logrus.Logger) *Service {
	return &Service{
		usersService: usersService,
		dbCollection: db.Collection("tasks"),
		log:          log,
	}
}

func (s *Service) CreateTask(ctx context.Context, task Task) (*Task, error) {
	log := s.log.WithContext(ctx).WithField("task", task)
	task.ID = primitive.NewObjectID().Hex()
	task.TimeAdded = time.Now()
	_, err := s.dbCollection.InsertOne(ctx, task)
	if err != nil {
		log.WithError(err).Error("failed to save task to db")
		return nil, err
	}
	return &task, nil
}

func (s *Service) GetTaskWithinTimeRange(ctx context.Context, userID string, startTime, endTime time.Time) (*Task, error) {
	filter := bson.M{"$or": []bson.M{
		{
			"startTime": bson.M{"$lte": startTime},
			"endTime":   bson.M{"$gte": startTime},
			"userId":    userID,
		},
		{
			"startTime": bson.M{"$lte": endTime},
			"endTime":   bson.M{"$gte": endTime},
			"userId":    userID,
		},
	}}
	var task Task
	err := s.dbCollection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		s.log.WithError(err).Error("cannot retrieve task from db")
		return nil, err
	}
	return &task, nil
}

func (s *Service) GetTask(ctx context.Context, taskID string) (*Task, error) {
	var task Task
	log := s.log.WithContext(ctx).WithField("taskId", taskID)
	err := s.dbCollection.FindOne(ctx, bson.M{"_id": taskID}).Decode(&task)
	if err != nil {
		log.WithError(err).Error("failed to retrieve task from db by id")
		return nil, err
	}
	return &task, nil
}

func (s *Service) GetTasks(ctx context.Context, userID, lastID string, limit int) ([]Task, error) {
	log := s.log.WithContext(ctx).WithField("userId", userID).WithField("lastId", lastID).
		WithField("limit", limit)
	_, err := s.usersService.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": bson.M{"$gt": lastID}, "userId": userID}
	findOpt := options.Find().SetLimit(int64(limit))
	cursor, err := s.dbCollection.Find(ctx, filter, findOpt)
	if err != nil {
		log.WithError(err).Error("failed for retrieve tasks from db")
		return nil, err
	}
	defer cursor.Close(ctx)
	var tasks []Task
	err = cursor.All(ctx, &tasks)
	if err != nil {
		log.WithError(err).Error("failed to decode retrieved tasks")
		return nil, err
	}
	return tasks, nil
}

func (s *Service) DeleteTask(ctx context.Context, taskID string) (*Task, error) {
	var task Task
	log := s.log.WithContext(ctx).WithField("taskId", taskID)
	err := s.dbCollection.FindOneAndDelete(ctx, bson.M{"_id": taskID}).Decode(&task)
	if err != nil {
		log.WithError(err).Error("failed to delete task from db")
		return nil, err
	}
	return &task, nil
}

func (s *Service) UpdateTask(ctx context.Context, taskID string, update Task) (*Task, error) {
	log := s.log.WithContext(ctx).WithField("taskId", taskID).WithField("update", update)
	filter := bson.M{"_id": taskID}
	updateBSON := bson.M{"$set": update}
	_, err := s.dbCollection.UpdateOne(ctx, filter, updateBSON)
	if err != nil {
		log.WithError(err).Error("failed to update task in db")
		return nil, err
	}
	return s.GetTask(ctx, taskID)
}
