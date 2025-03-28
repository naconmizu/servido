package crates

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// connectMongoDB establishes a connection to MongoDB Atlas.
func connectMongoDB(uri string) (*mongo.Client, context.Context, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}
	return client, context.Background(), nil // Fresh context for long-lived operations
}
