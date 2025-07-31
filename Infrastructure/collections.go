package infrastructure

import (
    "go.mongodb.org/mongo-driver/mongo"
)

func UserCollection() *mongo.Collection {
    return MongoDB.Collection("users")
}