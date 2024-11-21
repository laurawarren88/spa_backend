package models

import (
	"context"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username" binding:"required"`
	Email     string             `json:"email" bson:"email" binding:"required"`
	Password  string             `json:"password" bson:"password" binding:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (u *User) Validate(ctx context.Context, db *mongo.Collection) map[string]string {
	errors := make(map[string]string)

	if u.Username == "" {
		errors["username"] = "Username is required"
	} else {
		u.Username = strings.ToLower(strings.TrimSpace(u.Username))
		if len(u.Username) < 3 || len(u.Username) > 100 {
			errors["username"] = "Username must be between 3 and 100 characters"
		}
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9_.-]+$`, u.Username)
		if !matched {
			errors["username"] = "Username can only contain letters, numbers, underscores, dashes, and periods"
		} else {
			// Check if username exists in the database
			filter := bson.M{"username": u.Username}
			var existingUser User
			err := db.FindOne(ctx, filter).Decode(&existingUser)
			if err == nil {
				errors["username"] = "Username already exists"
			}
		}
	}

	if u.Email == "" {
		errors["email"] = "Email is required"
	} else {
		u.Email = strings.ToLower(strings.TrimSpace(u.Email))
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(emailRegex, u.Email)
		if !matched {
			errors["email"] = "Invalid email format"
		} else {
			// Check if email exists in the database
			filter := bson.M{"email": u.Email}
			var existingUser User
			err := db.FindOne(ctx, filter).Decode(&existingUser)
			if err == nil {
				errors["email"] = "Email already exists"
			}
		}
	}

	if u.Password == "" {
		errors["password"] = "Password is required"
	} else {
		if len(u.Password) < 6 {
			errors["password"] = "Password must be at least 6 characters long"
		}
		passwordRegex := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).+$`
		matched, _ := regexp.MatchString(passwordRegex, u.Password)
		if !matched {
			errors["password"] = "Password must contain at least one lowercase letter, one uppercase letter, and one number"
		}
	}

	return errors
}
