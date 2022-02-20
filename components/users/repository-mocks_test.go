package users_test

import (
	"context"

	"github.com/wisdommatt/creativeadvtech-assessment/components/users"
)

// userRepoMock is the user repository mock for mocking the user
// repository.
// this mock is created manually for simplicity sake, in a real
// application the mocks will/should be generated automatically from
// interfaces using a tool like https://github.com/vektra/mockery
type userRepoMock struct {
	saveUserFunc       func(ctx context.Context, user users.User) (*users.User, error)
	getUserByIDFunc    func(ctx context.Context, userID string) (*users.User, error)
	getUsersFunc       func(ctx context.Context, lastID string, limit int) ([]users.User, error)
	deleteUserByIDFunc func(ctx context.Context, userID string) (*users.User, error)
}

func (r *userRepoMock) saveUser(ctx context.Context, user users.User) (*users.User, error) {
	return r.saveUserFunc(ctx, user)
}

func (r *userRepoMock) getUserByID(ctx context.Context, userID string) (*users.User, error) {
	return r.getUserByIDFunc(ctx, userID)
}

func (r *userRepoMock) getUsers(ctx context.Context, lastID string, limit int) ([]users.User, error) {
	return r.getUsersFunc(ctx, lastID, limit)
}

func (r *userRepoMock) deleteUserByID(ctx context.Context, userID string) (*users.User, error) {
	return r.deleteUserByIDFunc(ctx, userID)
}
