package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	BookID primitive.ObjectID `json:"book_id" bson:"book_id"` // Add this
	// UserID    primitive.ObjectID `json:"user_id" bson:"user_id"` // Add this
	Review    string    `json:"review" bson:"review"`
	Rating    int       `json:"rating" bson:"rating" binding:"required,min=1,max=5"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (r *Review) Validate() map[string]string {
	errors := make(map[string]string)

	if r.Review == "" {
		errors["review"] = "Review is required"
	}

	if r.Rating == 0 {
		errors["rating"] = "Rating is required"
	}

	return errors
}
