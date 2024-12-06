package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Username  string             `json:"username" bson:"username"`
	Review    string             `json:"review" bson:"review"`
	Rating    int                `json:"rating" bson:"rating" binding:"required,min=1,max=5"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
	Book      Book               `json:"book" bson:"book"`
	User      User               `json:"user" bson:"user"`
}

func (r *Review) Validate() map[string]string {
	errors := make(map[string]string)
	if r.Review == "" {
		errors["review"] = "Review is required"
	}
	if r.Rating < 1 || r.Rating > 5 {
		errors["rating"] = "Rating must be between 1 and 5"
	}
	return errors
}
