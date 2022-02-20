package users

import (
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// Service is the interface that describes a user
// service object.
type Service interface {
}

// userService is the default implementation for Service
// interface.
type userService struct {
	userRepo repository
	log      *logrus.Logger
}

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
