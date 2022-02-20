package users

import "go.mongodb.org/mongo-driver/mongo"

// types, variables and functions related to user repository
// are not exported because they are not supposed to be exposed
// outside of this package.

// repository is the interface that describes a user
// repository.
type repository interface {
}

// userRepo is the default implementation for repository
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
