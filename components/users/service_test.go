package users

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
)

// tests for other methods in the user service will follow
// this same pattern.
func Test_userService_CreateUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		user User
	}
	tests := []struct {
		name     string
		args     args
		want     *User
		wantErr  bool
		userRepo *userRepoMock
	}{
		{
			name: "saveUser implementation with error",
			args: args{ctx: context.TODO(), user: User{
				FirstName: "Hello",
				LastName:  "World",
			}},
			wantErr: true,
			userRepo: &userRepoMock{
				saveUserFunc: func(ctx context.Context, user User) (*User, error) {
					return nil, fmt.Errorf("invalid user")
				},
			},
		},
		{
			name: "saveUser implementation without error",
			args: args{ctx: context.TODO(), user: User{
				FirstName: "Hello",
				LastName:  "World",
			}},
			userRepo: &userRepoMock{
				saveUserFunc: func(ctx context.Context, user User) (*User, error) {
					user.ID = "user-id"
					return &user, nil
				},
			},
			want: &User{
				ID:        "user-id",
				FirstName: "Hello",
				LastName:  "World",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{userRepo: tt.userRepo, log: logrus.New()}
			got, err := s.CreateUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userService.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
