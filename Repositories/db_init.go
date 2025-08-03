package repositories

import (
	infrastructure "blog-api/Infrastructure"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient() *mongo.Client {
	uri := infrastructure.Env.MONGODB_URI
	if uri == "" {
		log.Fatal("MONGODB_URI not set in environment")
	}

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal("MongoDB connection error: ", err)
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal("MongoDB ping error: ", err)
	}

	return client
}
