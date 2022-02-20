package users

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
// types, variables and functions related to user repository
// are not exported because they are not supposed to be exposed
// outside of this package.

// repository is the interface that describes a user
// repository.
type repository interface {
	saveUser(ctx context.Context, user User) (*User, error)
	getUserByID(ctx context.Context, userID string) (*User, error)
	getUsers(ctx context.Context, lastID string, limit int) ([]User, error)
	deleteUserByID(ctx context.Context, userID string) (*User, error)
}

// userRepo is the default implementation for user repository
// interface.
type userRepo struct {
	usersCollection *mongo.Collection
}

// newRepository creates a new user repository.
func newRepository(db *mongo.Database) *userRepo {
	return &userRepo{
		usersCollection: db.Collection("users"),
	}
}

func (r *userRepo) saveUser(ctx context.Context, user User) (*User, error) {
	user.ID = primitive.NewObjectID().Hex()
	user.TimeAdded = time.Now()
	user.LastUpdated = time.Now()
	_, err := r.usersCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) getUserByID(ctx context.Context, userID string) (*User, error) {
	filter := bson.M{"_id": userID}
	var user User
	err := r.usersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) getUsers(ctx context.Context, lastID string, limit int) ([]User, error) {
	filter := bson.M{"_id": bson.M{"$gt": lastID}}
	findOpt := options.Find().SetLimit(int64(limit))
	cursor, err := r.usersCollection.Find(ctx, filter, findOpt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var users []User
	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// deleteUserByID deletes a single user from the database by id.
//
// in the case of production applications there might be some business
// requirements/needs that the user details should not be deleted fully from
// the database, in that case we can either have a separate trash database / collection
// for storing deleted users.
func (r *userRepo) deleteUserByID(ctx context.Context, userID string) (*User, error) {
	filter := bson.M{"_id": userID}
	cursor := r.usersCollection.FindOneAndDelete(ctx, filter)
	if cursor.Err() != nil {
		return nil, cursor.Err()
	}
	var deletedUser User
	err := cursor.Decode(&deletedUser)
	if err != nil {
		return nil, err
	}
	return &deletedUser, nil
}
