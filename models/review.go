package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Review    string             `json:"review" bson:"review" binding:"required"`
	Rating    int                `json:"rating" bson:"rating" binding:"required,min=1,max=5"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
	Book      Book               `json:"book" bson:"book"`
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
