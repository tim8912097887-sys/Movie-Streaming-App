package users

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type repository struct{
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *repository {
	return &repository{
		collection: collection,
	}
}

func (r *repository) CreateUser(ctx context.Context, user CreateUserSchema) (User, error) {
    result, err :=r.collection.InsertOne(ctx, user)
	
	if err != nil {
		return User{}, err
	}

	id, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return User{}, errors.New("failed to convert inserted id")
	}

	return User{
		ID:              id,
		Name:            user.Name,
		Email:           user.Email,
	}, nil
}

func (r *repository) FindUserByEmail(ctx context.Context, email string) (User, error) {
	
	var user User

	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	return user, nil
}

func (r *repository) FindUserById(ctx context.Context, id string) (User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return User{}, err
	}
	var user User

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}
	return user, nil
}

func (r *repository) UpdateUser(ctx context.Context, user UpdateUserSchema) error {
	objectID, err := bson.ObjectIDFromHex(user.Id)
	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{"token_version": user.TokenVersion}}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	return nil
}