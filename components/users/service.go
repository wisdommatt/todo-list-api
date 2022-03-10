package users

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Service is the interface that describes a user
// service object.
type Service interface {
	CreateUser(ctx context.Context, user User) (*User, error)
	GetUser(ctx context.Context, userID string) (*User, error)
	GetUsers(ctx context.Context, lastID string, limit int) ([]User, error)
	DeleteUser(ctx context.Context, userID string) (*User, error)
}

// userService is the default implementation for Service
// interface.
type userService struct {
	userRepo repository
	log      *logrus.Logger
}

var (
	errSomethingWentWrong = fmt.Errorf("an error occured, please try again later")
)

// NewService creates a new user service.
//
// an instance of mongodb is need when creating a new user
// service because the user service depends on the user repository
// which in this case depends on mongodb.
// the user service does not interact directly with mongodb in anyway.
func NewService(mongoDB *mongo.Database, log *logrus.Logger) *userService {
	return &userService{
		userRepo: newRepository(mongoDB),
		log:      log,
	}
}

func (s *userService) CreateUser(ctx context.Context, user User) (*User, error) {
	log := s.log.WithContext(ctx).WithField("user", user)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error()
		return nil, errSomethingWentWrong
	}
	user.Password = string(hashedPassword)
	newUser, err := s.userRepo.saveUser(ctx, user)
	if err != nil {
		log.WithError(err).Error("an error occured while creating new user")
		return nil, errSomethingWentWrong
	}
	return newUser, nil
}

func (s *userService) GetUser(ctx context.Context, userID string) (*User, error) {
	log := s.log.WithContext(ctx).WithField("userId", userID)
	user, err := s.userRepo.getUserByID(ctx, userID)
	if err != nil {
		log.WithError(err).Error("an error occured while retrieving user")
		return nil, fmt.Errorf("user does not exist")
	}
	return user, nil
}

func (s *userService) GetUsers(ctx context.Context, lastID string, limit int) ([]User, error) {
	log := s.log.WithContext(ctx).WithField("lastId", lastID).WithField("limit", limit)
	users, err := s.userRepo.getUsers(ctx, lastID, limit)
	if err != nil {
		log.WithError(err).Error("an error occured while retrieving users")
		return nil, errSomethingWentWrong
	}
	return users, nil
}

func (s *userService) DeleteUser(ctx context.Context, userID string) (*User, error) {
	log := s.log.WithContext(ctx).WithField("userId", userID)
	// checking if user id is valid.
	_, err := s.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.deleteUserByID(ctx, userID)
	if err != nil {
		log.WithError(err).Error("an error occured while deleting user")
		return nil, errSomethingWentWrong
	}
	return user, nil
}
