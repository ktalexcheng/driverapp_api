package util

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Struct for MongoDB client
type MongoClient struct {
	// Client   *mongo.Client
	MongoDB         *mongo.Database
	MongoURI        string
	Database        string
	RideDataColl    *mongo.Collection
	RideRecordsColl *mongo.Collection
	UsersColl       *mongo.Collection
}

// Initializes a new MongoDB connection client
func NewMongoClient(connStr string, dbName string, certPath string) (*MongoClient, error) {
	credential := options.Credential{
		AuthMechanism: "MONGODB-X509",
		AuthSource:    "external",
	}
	dbUri := fmt.Sprintf(
		"%s/?tlsCertificateKeyFile=%s",
		connStr,
		certPath,
	)
	clientOptions := options.Client().ApplyURI(dbUri).SetAuth(credential)

	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		panic(err)
	}
	log.Info().Msg("Succesfully connected to and pinged MongoDB")

	mongoClient := MongoClient{
		// Client:   client,
		MongoDB:         client.Database(dbName),
		MongoURI:        connStr,
		Database:        dbName,
		RideDataColl:    client.Database(dbName).Collection("rideData"),
		RideRecordsColl: client.Database(dbName).Collection("rideRecords"),
		UsersColl:       client.Database(dbName).Collection("users"),
	}

	return &mongoClient, nil
}
