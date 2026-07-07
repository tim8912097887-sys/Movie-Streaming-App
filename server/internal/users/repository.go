package users

import (
	"context"
	"errors"

	"github.com/tim8912097887-sys/server/internal/shared"
	"github.com/tim8912097887-sys/server/internal/shared/types"
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

func (r *repository) CreateUser(ctx context.Context, user types.CreateUserSchema) (types.User, error) {
    result, err :=r.collection.InsertOne(ctx, user)
	
	if err != nil {
		return types.User{}, err
	}

	id, ok := result.InsertedID.(bson.ObjectID)
	if !ok {
		return types.User{}, errors.New("failed to convert inserted id")
	}

	return types.User{
		ID:              id,
		Name:            user.Name,
		Email:           user.Email,
	}, nil
}

func (r *repository) FindUserByEmail(ctx context.Context, email string) (types.User, error) {
	
	var user types.User

	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return types.User{}, shared.ErrUserNotFound
		}
		return types.User{}, err
	}

	return user, nil
}

func (r *repository) FindUserById(ctx context.Context, id string) (types.User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return types.User{}, err
	}
	var user types.User

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return types.User{}, shared.ErrUserNotFound
		}
		return types.User{}, err
	}
	return user, nil
}

func (r *repository) UpdateUser(ctx context.Context, user types.UpdateUserSchema) error {
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