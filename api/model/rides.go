package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Struct for MongoDB rideRecords document mapping
type RideRecord struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	RideName  string             `bson:"rideName" json:"rideName"`
	RideDate  time.Time          `bson:"rideDate" json:"rideDate"`
	RideMeta  RideMeta           `bson:"rideMeta,OmitEmpty" json:"rideMeta,omitempty"`
	RideScore RideScore          `bson:"rideScore,OmitEmpty" json:"rideScore,omitempty"`
}

// Struct for MongoDB rideData document mapping
type RideDatum struct {
	Timestamp      time.Time          `bson:"timestamp" json:"timestamp"`
	RideRecordID   primitive.ObjectID `bson:"rideRecordId" json:"rideRecordId"`
	LocationLat    float64            `bson:"locationLat" json:"locationLat"`
	LocationLong   float64            `bson:"locationLong" json:"locationLong"`
	AccelerometerX float64            `bson:"accelerometerX" json:"accelerometerX"`
	AccelerometerY float64            `bson:"accelerometerY" json:"accelerometerY"`
	AccelerometerZ float64            `bson:"accelerometerZ" json:"accelerometerZ"`
	GyroscopeX     float64            `bson:"gyroscopeX" json:"gyroscopeX"`
	GyroscopeY     float64            `bson:"gyroscopeY" json:"gyroscopeY"`
	GyroscopeZ     float64            `bson:"gyroscopeZ" json:"gyroscopeZ"`
	RotationX      float64            `bson:"rotationX" json:"rotationX"`
	RotationY      float64            `bson:"rotationY" json:"rotationY"`
	RotationZ      float64            `bson:"rotationZ" json:"rotationZ"`
	RotationW      float64            `bson:"rotationW" json:"rotationW"`
}

// Struct for rideMeta field in rideRecords documents
type RideMeta struct {
	Distance        float64 `bson:"distance" json:"distance"`
	Duration        float64 `bson:"duration" json:"duration"`
	MaxAcceleration float64 `bson:"maxAcceleration" json:"maxAcceleration"`
	AccelerationRMS float64 `bson:"accelerationRms" json:"accelerationRms"`
}

// Struct for rideScore field in rideRecords documents
type RideScore struct {
	Overall      float64 `bson:"overall" json:"overall"`
	Speed        float64 `bson:"speed" json:"speed"`
	Comfort      float64 `bson:"comfort" json:"comfort"`
	Acceleration float64 `bson:"acceleration" json:"acceleration"`
	Braking      float64 `bson:"braking" json:"braking"`
	Cornering    float64 `bson:"cornering" json:"cornering"`
}

// Struct for JSON response at GET /rides/{rideId} endpoint
type Ride struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	RideName  string             `bson:"rideName" json:"rideName"`
	RideDate  time.Time          `bson:"rideDate" json:"rideDate"`
	RideMeta  RideMeta           `bson:"rideMeta" json:"rideMeta"`
	RideScore RideScore          `bson:"rideScore" json:"rideScore"`
	RideData  []RideDatum        `bson:"rideData" json:"rideData"`
}

// Struct for Judge service response
type JudgeResult struct {
	RideScore RideScore `json:"rideScore"`
	RideMeta  RideMeta  `json:"rideMeta"`
}
