package infrastructure

import (
    "context"
    "fmt"
    "log"
    "time"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func InitMongo() {
    mongoURI := Env.MONGODB_URI
    dbName := Env.DB_NAME

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatalf("MongoDB connection error: %v", err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatalf("MongoDB ping error: %v", err)
    }

    MongoClient = client
    MongoDB = client.Database(dbName)

    fmt.Println("Connected to MongoDB")
}
