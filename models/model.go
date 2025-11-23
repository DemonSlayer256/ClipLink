package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JWT struct {
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
	jwt.RegisteredClaims
}

type User struct {
	Id   primitive.ObjectID `bson:"_id" json:"id"`
	User string             `bson:"user" json:"user"`
	Pass string             `bson:"pass" json:"pass"`
	Left int                `bson:"left" json:"left"`
}

type URLMapping struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"`
	Link       string             `bson:"link" json:"link"`
	Shorted    string             `bson:"shorted" json:"shorted"`
	Expires_at time.Time          `bson:"expiresAt,omitempty" json:"expiresAt,omitempty"`
	Created_at time.Time          `bson:"createAt" json:"createdAt"`
	User       string             `bson:"user" json:"user"`
}
