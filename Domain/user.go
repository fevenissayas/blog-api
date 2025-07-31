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
    ID           primitive.ObjectID
    Username     string
    Email        string
    Password     string
    Role         Role
    CreatedAt    time.Time
    UpdatedAt    time.Time
    Bio          string
    ProfilePicture string
	  ContactInfo  string
}
