package infrastructure

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func InitMongo() {
    _ = godotenv.Load()
    mongoURI := os.Getenv("MONGODB_URI")
    dbName := os.Getenv("DB_NAME")

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
