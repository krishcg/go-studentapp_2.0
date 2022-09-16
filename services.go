package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func add(x int, y int) int {
	return x + y
}

func MongoDBConnection(clientOptions *options.ClientOptions) *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, clientOptions)
	return client
}
