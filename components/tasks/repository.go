package tasks

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// the repository layer is where all interactions with any form of
// persistent data storage/retrieval happens.
// types, variables and functions related to task repository
// are not exported because they are not supposed to be exposed
// outside of this package.

// repository is the interface that describes a task
// repository.
type repository interface {
	saveTask(ctx context.Context, task Task) (*Task, error)
	getTaskByID(ctx context.Context, taskID string) (*Task, error)
	getTasks(ctx context.Context, userID, lastID string, limit int) ([]Task, error)
	deleteTaskByID(ctx context.Context, taskID string) (*Task, error)
	getTaskWithTimeRange(ctx context.Context, userID string, startTime, endTime time.Time) (*Task, error)
	updateTaskByID(ctx context.Context, taskID string, update Task) error
}

// taskRepo is the default implementation for task repository
// interface.
type taskRepo struct {
	tasksCollection *mongo.Collection
}

// newRepository creates a new task repository.
func newRepository(db *mongo.Database) *taskRepo {
	return &taskRepo{
		tasksCollection: db.Collection("tasks"),
	}
}

func (r *taskRepo) saveTask(ctx context.Context, task Task) (*Task, error) {
	task.ID = primitive.NewObjectID().Hex()
	task.TimeAdded = time.Now()
	_, err := r.tasksCollection.InsertOne(ctx, task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepo) getTaskByID(ctx context.Context, taskID string) (*Task, error) {
	filter := bson.M{"_id": taskID}
	var task Task
	err := r.tasksCollection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepo) getTasks(ctx context.Context, userID, lastID string, limit int) ([]Task, error) {
	filter := bson.M{"_id": bson.M{"$gt": lastID}, "userId": userID}
	findOpt := options.Find().SetLimit(int64(limit))
	cursor, err := r.tasksCollection.Find(ctx, filter, findOpt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var tasks []Task
	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepo) deleteTaskByID(ctx context.Context, taskID string) (*Task, error) {
	filter := bson.M{"_id": taskID}
	var task Task
	err := r.tasksCollection.FindOneAndDelete(ctx, filter).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepo) getTaskWithTimeRange(ctx context.Context, userID string, startTime, endTime time.Time) (*Task, error) {
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
	err := r.tasksCollection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepo) updateTaskByID(ctx context.Context, taskID string, update Task) error {
	filter := bson.M{"_id": taskID}
	updateBSON := bson.M{"$set": update}
	_, err := r.tasksCollection.UpdateOne(ctx, filter, updateBSON)
	return err
}
