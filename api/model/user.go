package model

import (
	"context"

	"github.com/ktalexcheng/trailbrake_api/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Struct for MongoDB users document mapping
type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserAlias string             `bson:"userAlias" json:"userAlias"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"password"`
}

type UserStats struct {
	RidesCount      int     `bson:"ridesCount" json:"ridesCount"`
	TotalDistance   float64 `bson:"totalDistance" json:"totalDistance"`
	TotalRideTime   float64 `bson:"totalRideTime" json:"totalRideTime"`
	MaxAcceleration float64 `bson:"maxAcceleration" json:"maxAcceleration"`
}

func (u *User) CheckUserExists(mg *util.MongoClient) (bool, error) {
	usersCol := mg.MongoDB.Collection("users")
	filter := bson.M{
		"$or": []bson.M{
			{"_id": u.ID},
			{"userId": u.UserAlias},
			{"email": u.Email},
		},
	}

	count, err := usersCol.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (u *User) ValidateUserPass(mg *util.MongoClient) (bool, error) {
	usersCol := mg.MongoDB.Collection("users")

	var userCred User
	err := usersCol.FindOne(context.TODO(), bson.D{{Key: "email", Value: u.Email}}).Decode(&userCred)
	if err != nil {
		return false, err
	}

	if u.Password == userCred.Password {
		// Only fetch user identity after credentials are validated
		u.ID = userCred.ID
		u.UserAlias = userCred.UserAlias

		return true, nil
	}

	return false, nil
}
