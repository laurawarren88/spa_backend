package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"password" bson:"password"`
	Role      string             `json:"role" bson:"role"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (u *User) Validate() map[string]string {
	errors := make(map[string]string)

	if u.Username == "" {
		errors["username"] = "Username is required"
	}
	if u.Email == "" {
		errors["email"] = "Email is required"
	}
	if u.Password == "" {
		errors["password"] = "password is required"
	}

	return errors
}
