package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Author      string             `json:"author" bson:"author"`
	Category    string             `json:"category" bson:"category"`
	Description string             `json:"description" bson:"description"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

func (b *Book) Validate() map[string]string {
	errors := make(map[string]string)

	if b.Title == "" {
		errors["title"] = "Title is required"
	}
	if b.Author == "" {
		errors["author"] = "Author is required"
	}
	if b.Category == "" {
		errors["category"] = "Category is required"
	}

	if b.Description == "" {
		errors["description"] = "Description is required"
	}

	return errors
}
