package user

import (
	"context"
	"go-gin-e-comm/common"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(c context.Context, email, password string) error
	FindByEmail(c context.Context, email string) (*User, error)
	CreateUserData(c context.Context, user *User, userID primitive.ObjectID) error
	FindUserByID(c context.Context, userID primitive.ObjectID) (*User, error)
}

type userRepository struct {
	db  *mongo.Database
	log *slog.Logger
}

func NewUserRepository(db *mongo.Database, log *slog.Logger) UserRepository {
	return &userRepository{db: db, log: log}
}

func (r *userRepository) FindUserByID(c context.Context, userID primitive.ObjectID) (*User, error) {
	var userDB User
	err := r.db.Collection("users").FindOne(c, bson.M{"_id": userID}).Decode(&userDB)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		} else {
			r.log.Error("Database error", "error", err.Error())
			return nil, common.ErrDatabase
		}
	}

	return &userDB, nil
}

func (r *userRepository) CreateUser(c context.Context, email, password string) error {
	user := User{
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := r.db.Collection("users").InsertOne(c, user)
	if err != nil {
		r.log.Error(
			"Database error when creating user",
			"email", email,
			"error", err.Error(),
		)
		return err
	}
	return nil
}

func (r *userRepository) CreateUserData(c context.Context, user *User, userID primitive.ObjectID) error {
	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"address":    user.Address,
			"updated_at": time.Now(),
		},
	}
	result, err := r.db.Collection("users").UpdateOne(c, filter, update)
	if err != nil {
		r.log.Error(
			"Failed when updating user data",
			"userID", userID.Hex(),
			"error", err.Error(),
		)
		return common.ErrDatabase
	}
	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) FindByEmail(c context.Context, email string) (*User, error) {
	var userDB User
	filter := bson.M{"email": email}
	err := r.db.Collection("users").FindOne(c, filter).Decode(&userDB)
	if err == mongo.ErrNoDocuments {
		r.log.Info("User by email not found", "email", email)
		return nil, nil
	}
	if err != nil {
		r.log.Error(
			"Database error when finding user by email",
			"email", email,
			"error", err.Error(),
		)
		return nil, err
	}
	return &userDB, nil
}
