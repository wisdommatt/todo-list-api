package users

import (
	"context"
)

// userRepoMock is the user repository mock for mocking the user
// repository.
// this mock is created manually for simplicity sake, in a real
// application the mocks will/should be generated automatically from
// interfaces using a tool like https://github.com/vektra/mockery
type userRepoMock struct {
	saveUserFunc       func(ctx context.Context, user User) (*User, error)
	getUserByIDFunc    func(ctx context.Context, userID string) (*User, error)
	getUserByEmailFunc func(ctx context.Context, email string) (*User, error)
	getUsersFunc       func(ctx context.Context, lastID string, limit int) ([]User, error)
	deleteUserByIDFunc func(ctx context.Context, userID string) (*User, error)
}

func (r *userRepoMock) saveUser(ctx context.Context, user User) (*User, error) {
	return r.saveUserFunc(ctx, user)
}

func (r *userRepoMock) getUserByID(ctx context.Context, userID string) (*User, error) {
	return r.getUserByIDFunc(ctx, userID)
}

func (r *userRepoMock) getUsers(ctx context.Context, lastID string, limit int) ([]User, error) {
	return r.getUsersFunc(ctx, lastID, limit)
}

func (r *userRepoMock) deleteUserByID(ctx context.Context, userID string) (*User, error) {
	return r.deleteUserByIDFunc(ctx, userID)
}

func (r *userRepoMock) getUserByEmail(ctx context.Context, email string) (*User, error) {
	return r.getUserByEmailFunc(ctx, email)
}
