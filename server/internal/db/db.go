package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)


func ConnectDB(ctx context.Context,dbUrl string) (*mongo.Client, error) {

	mongoClient, err := mongo.Connect(options.Client().ApplyURI(dbUrl))

	if err != nil {
		return nil, err
	}


	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		// Clean up connection after error
		if err = mongoClient.Disconnect(ctx); err != nil {
			return nil, err
		}
		return nil, err
	}
	return mongoClient, nil
}