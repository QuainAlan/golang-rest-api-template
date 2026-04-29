package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func mongoConnectionURI() string {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return "mongodb://localhost:27017"
	}
	return uri
}

// SetupMongoDB connects to MongoDB and returns the application logs collection.
// It reads MONGODB_URI (e.g. mongodb://mongo:27017 under Docker Compose); if unset,
// it defaults to mongodb://localhost:27017 for local runs against a host-mapped port.
func SetupMongoDB() *mongo.Collection {
	clientOptions := options.Client().ApplyURI(mongoConnectionURI())
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	return client.Database("logging").Collection("logs")
}
