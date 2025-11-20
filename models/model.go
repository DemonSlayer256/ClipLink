package models

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type JWT struct {
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
	jwt.RegisteredClaims
}

type User struct {
	Id   bson.ObjectID `bson:"_id"`
	User string        `bson:"user"`
	Pass string        `bson:"pass"`
	Left int           `bson:"left"`
}

type URLMapping struct {
	Id         bson.ObjectID `bson:"_id"`
	Link       string        `bson:"link"`
	Shorted    string        `bson:"shorted"`
	Expires_at time.Time     `bson:"expiresAt, omitempty"`
	Created_at time.Time     `bson:"createAt"`
	User       string        `bson:"user"`
}
