package domain

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type Role string

const (
    RoleUser  Role = "user"
    RoleAdmin Role = "admin"
)

type User struct {
    ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Username     string             `bson:"username" json:"username"`
    Email        string             `bson:"email" json:"email"`
    Password     string             `bson:"password,omitempty" json:"-"`
    Role         Role               `bson:"role" json:"role"`
    CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
    Bio          string             `bson:"bio,omitempty" json:"bio,omitempty"`
    ProfileImage string             `bson:"profile_image,omitempty" json:"profile_image,omitempty"`
}
