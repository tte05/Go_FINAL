package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Transaction struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Project  string             `json:"project"`
	Customer Customer           `json:"customer"`
	Status   string             `json:"status"`
	Date     time.Time          `json:"date"`
}

type Item struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type Customer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
