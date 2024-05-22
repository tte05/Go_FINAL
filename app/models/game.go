package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Game struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Genre       string             `bson:"genre" json:"genre"`
	Rating      float64            `bson:"rating,omitempty" json:"rating,omitempty"`
	Developer   string             `bson:"developer,omitempty" json:"developer,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
}
