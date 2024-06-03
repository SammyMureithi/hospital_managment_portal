package models

import (
	"time"
)

// User struct to define a user model
type Doctor struct {
    ID        string    `json:"id"` 
    FirstName string    `json:"first_name" bson:"first_name" validate:"required,min=3,max=20"`
    LastName  string    `json:"last_name" bson:"last_name" validate:"required"`
    Email     string    `json:"email" bson:"email" validate:"required,email"`
	Phone     string    `json:"phone" bson:"phone" validate:"required,len=10"`
    Available bool      `bson:"available"`
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
