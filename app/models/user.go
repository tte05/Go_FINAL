package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Username          string             `bson:"username" json:"username"`
	Password          string             `bson:"password,omitempty" json:"-"`
	Email             string             `bson:"email" json:"email"`
	Role              string             `bson:"role" json:"role"`
	Confirmed         bool               `bson:"confirmed" json:"confirmed"`
	ConfirmationToken string             `bson:"confirmation_token" json:"-"`
	ResetToken        string             `bson:"reset_token,omitempty" json:"-"`
	ResetTokenExpiry  time.Time          `bson:"reset_token_expiry,omitempty" json:"-"`
}
