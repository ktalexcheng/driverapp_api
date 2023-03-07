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
	MongoDB  *mongo.Database
	MongoURI string
	Database string
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
		MongoDB:  client.Database(dbName),
		MongoURI: connStr,
		Database: dbName,
	}

	return &mongoClient, nil
}
