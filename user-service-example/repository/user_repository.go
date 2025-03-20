package repository

import (
	"context"

	"github.com/zuyatna/edu-connect/user-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserRepository interface {
	RegisterUser(ctx context.Context, user *model.User) (*model.User, error)
	LoginUser(ctx context.Context, email, password string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateOrderCountUser(ctx context.Context, id string) error
	DeleteUser(ctx context.Context, id string) error
}

type UserRepository struct {
	usersCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		usersCollection: db.Collection("users"),
	}
}

func (u *UserRepository) RegisterUser(ctx context.Context, user *model.User) (*model.User, error) {
	doc := bson.D{
		{Key: "name", Value: user.Name},
		{Key: "email", Value: user.Email},
		{Key: "password", Value: user.Password},
		{Key: "order_count", Value: 0},
	}

	result, err := u.usersCollection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return user, nil
}

func (u *UserRepository) LoginUser(ctx context.Context, email, password string) (*model.User, error) {
	var user model.User

	err := u.usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	err = user.CompareHashAndPassword(password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = u.usersCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := u.usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	doc := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: user.Name},
			{Key: "email", Value: user.Email},
			{Key: "password", Value: user.Password},
		}},
	}

	_, err := u.usersCollection.UpdateOne(ctx, bson.M{"_id": user.ID}, doc)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) UpdateOrderCountUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	doc := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "order_count", Value: 1},
		}},
	}

	_, err = u.usersCollection.UpdateOne(ctx, bson.M{"_id": objectID}, doc)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) DeleteUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = u.usersCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	return nil
}
