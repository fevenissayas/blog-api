package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
    RoleUser  Role = "user"
    RoleAdmin Role = "admin"
)

type User struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Username        string             `bson:"username" json:"username" validate:"required,min=3,max=20"`
    Email           string             `bson:"email" json:"email" validate:"required,email"`
    Password        string             `bson:"password,omitempty" json:"password,omitempty" validate:"required,min=6"`
    Role            Role               `bson:"role" json:"role"`
    CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
    Bio             string             `bson:"bio,omitempty" json:"bio,omitempty"`
    ProfilePicture  string             `bson:"profile_picture,omitempty" json:"profile_picture,omitempty"`
    ContactInfo     string             `bson:"contact_info,omitempty" json:"contact_info,omitempty"`
}
