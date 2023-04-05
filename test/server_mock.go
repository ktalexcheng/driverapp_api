package test

import (
	"context"
	"net/http/httptest"

	"github.com/ktalexcheng/trailbrake_api/api/router"
	"github.com/ktalexcheng/trailbrake_api/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mim "github.com/ktalexcheng/dp-mongodb-in-memory"
)

const testDbName = "test"

func NewMockDB() *util.MongoClient {
	mimDb, err := mim.Start(context.TODO(), "6.0.0")
	if err != nil {
		panic("unable to start in-memory database")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mimDb.URI()))
	if err != nil {
		panic("unable to connect to in-memory database")
	}

	mongoClient := util.MongoClient{
		MongoDB:         client.Database(testDbName),
		MongoURI:        "",
		Database:        testDbName,
		RideDataColl:    client.Database(testDbName).Collection("rideData"),
		RideRecordsColl: client.Database(testDbName).Collection("rideRecords"),
		UsersColl:       client.Database(testDbName).Collection("users"),
	}

	return &mongoClient
}

func NewTestServer(mg *util.MongoClient) *httptest.Server {
	r := router.Router(mg)
	testServer := httptest.NewServer(r)

	return testServer
}
