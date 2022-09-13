package users

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wisdommatt/todo-list-api/internal/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	FirstName   string    `json:"firstName" bson:"firstName,omitempty"`
	LastName    string    `json:"lastName" bson:"lastName,omitempty"`
	Email       string    `json:"email" bson:"email,omitempty"`
	Password    string    `json:"-" bson:"password,omitempty"`
	TimeAdded   time.Time `json:"timeAdded" bson:"timeAdded,omitempty"`
	LastUpdated time.Time `json:"-" bson:"lastUpdated,omitempty"`
}

type Service struct {
	log          *logrus.Logger
	dbCollection *mongo.Collection
}

func NewUsersService(db *mongo.Database, log *logrus.Logger) *Service {
	return &Service{
		log:          log,
		dbCollection: db.Collection("users"),
	}
}

func (s *Service) CreateUser(ctx context.Context, user User) (*User, error) {
	log := s.log.WithContext(ctx).WithField("user", user)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("cannot generate password from hash")
		return nil, err
	}
	user.Password = string(hashedPassword)
	user.ID = primitive.NewObjectID().Hex()
	user.TimeAdded = time.Now()
	user.LastUpdated = time.Now()
	_, err = s.dbCollection.InsertOne(ctx, user)
	if err != nil {
		log.WithError(err).Error("cannot save user to db")
		return nil, err
	}
	return &user, nil
}

func (s *Service) GetUser(ctx context.Context, userID string) (*User, error) {
	var user User
	log := s.log.WithContext(ctx).WithField("userId", userID)
	err := s.dbCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		log.WithError(err).Error("cannot retrieve user from db by id")
		return nil, err
	}
	return &user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	log := s.log.WithContext(ctx).WithField("email", email)
	err := s.dbCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		log.WithError(err).Error("cannot retrieve user from db by email")
		return nil, err
	}
	return &user, nil
}

func (s *Service) GetUsers(ctx context.Context, lastID string, limit int) ([]User, error) {
	log := s.log.WithContext(ctx).WithField("lastId", lastID).WithField("limit", limit)
	filter := bson.M{"_id": bson.M{"$gt": lastID}}
	findOpt := options.Find().SetLimit(int64(limit))
	cursor, err := s.dbCollection.Find(ctx, filter, findOpt)
	if err != nil {
		log.WithError(err).Error("cannot retrieve users from db")
		return nil, err
	}
	defer cursor.Close(ctx)
	var users []User
	err = cursor.All(ctx, &users)
	if err != nil {
		log.WithError(err).Error("cannot decode retrieved users")
		return nil, err
	}
	return users, nil
}

func (s *Service) DeleteUser(ctx context.Context, userID string) (*User, error) {
	log := s.log.WithContext(ctx).WithField("userId", userID)
	filter := bson.M{"_id": userID}
	var deletedUser User
	err := s.dbCollection.FindOneAndDelete(ctx, filter).Decode(&deletedUser)
	if err != nil {
		log.WithError(err).Error("failed to delete user from db")
		return nil, err
	}
	return &deletedUser, nil
}

func (s *Service) LoginUser(ctx context.Context, email, password string) (*User, string, error) {
	var errInvalidCredentials = fmt.Errorf("invalid credentials")
	log := s.log.WithField("email", email)
	userWithEmail, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", errInvalidCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(userWithEmail.Password), []byte(password))
	if err != nil {
		log.WithError(err).Error("failed to compare password with hash")
		return nil, "", errInvalidCredentials
	}
	authToken, err := jwt.Encode([]byte(os.Getenv("JWT_SECRET")), jwt.Payload{UserID: userWithEmail.ID})
	if err != nil {
		log.WithError(err).Error("failed to encode jwt token")
		return nil, "", err
	}
	return userWithEmail, authToken, nil
}
