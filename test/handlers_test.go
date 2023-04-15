package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ktalexcheng/trailbrake_api/api/middleware"
	"github.com/ktalexcheng/trailbrake_api/api/model"
	"github.com/ktalexcheng/trailbrake_api/util"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var apiServer *httptest.Server
var testUser *model.User
var authHeader map[string]string

func sendRequestToMockServer(t *testing.T, mg *util.MongoClient, method string, endpoint string, body io.Reader, headers interface{}) *httptest.ResponseRecorder {
	fmt.Println(apiServer.URL)
	// Create new request
	req, err := http.NewRequest(method, apiServer.URL+endpoint, body)
	if err != nil {
		t.Errorf("Failed to create request: %v", err)
	}

	// Set headers
	if _, ok := headers.(map[string]string); ok {
		for k, v := range headers.(map[string]string) {
			req.Header.Set(k, v)
		}
	}

	// Create response recorder and send request
	rr := httptest.NewRecorder()
	apiServer.Config.Handler.ServeHTTP(rr, req)

	return rr
}

func unmarshalJSON(rr *httptest.ResponseRecorder, target any) error {
	err := json.NewDecoder(rr.Body).Decode(target)
	if err != nil {
		return err
	}

	return nil
}

func initTestEnv() {
	testEnvs := map[string]string{
		"MONGO_DB_NAME":    "driverAppDB",
		"TOKEN_SECRET_KEY": "ead1a39f200400e43f7f3da657b42f8a2243d67be6343ac4209b05636b9ad426",
		"JUDGE_URL":        "https://trailbrake-judge-f6muv3fwlq-de.a.run.app",
	}
	for k, v := range testEnvs {
		util.SetEnvIfMissing(k, v)
	}
}

func createTestUser(mg *util.MongoClient, user *model.User) error {
	_, err := mg.UsersColl.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	return nil
}

func createTestRide(mg *util.MongoClient, user *model.User) (*model.Ride, error) {
	testRideMeta := model.RideMeta{
		Distance:        100,
		Duration:        60,
		MaxAcceleration: 9.8,
	}
	testRideScore := model.RideScore{
		Overall:      88,
		Speed:        100,
		Acceleration: 85,
		Braking:      70,
		Cornering:    90,
	}

	rideId := primitive.NewObjectID()
	rideName := "Test ride"
	rideDate := time.Now()
	testRideRecord := model.RideRecord{
		ID:        rideId,
		UserID:    user.ID,
		RideName:  rideName,
		RideDate:  rideDate,
		RideMeta:  testRideMeta,
		RideScore: testRideScore,
	}
	testRideData := []model.RideDatum{
		{
			Timestamp:      rideDate.Add(time.Second),
			RideRecordID:   rideId,
			GyroscopeX:     0.001,
			GyroscopeY:     -0.002,
			GyroscopeZ:     0.003,
			AccelerometerX: 0.004,
			AccelerometerY: 0.005,
			AccelerometerZ: -0.006,
			LocationLat:    25.105497,
			LocationLong:   121.597366,
		},
		{
			Timestamp:      rideDate.Add(time.Second * 2),
			RideRecordID:   rideId,
			GyroscopeX:     0.002,
			GyroscopeY:     -0.004,
			GyroscopeZ:     0.006,
			AccelerometerX: 0.008,
			AccelerometerY: 0.010,
			AccelerometerZ: -0.012,
			LocationLat:    25.105597,
			LocationLong:   121.598366,
		},
	}

	_, err := mg.RideRecordsColl.InsertOne(context.TODO(), testRideRecord)
	if err != nil {
		return nil, err
	}

	var insertDocs []interface{}
	for _, r := range testRideData {
		insertDocs = append(insertDocs, r)
	}
	_, err = mg.RideDataColl.InsertMany(context.TODO(), insertDocs)
	if err != nil {
		return nil, err
	}

	testRide := &model.Ride{
		ID:        rideId,
		RideName:  rideName,
		RideDate:  rideDate,
		RideMeta:  testRideMeta,
		RideScore: testRideScore,
		RideData:  testRideData,
	}

	return testRide, nil
}

func TestMain(m *testing.M) {
	initTestEnv()

	// Prepare a fake user to pass the authentication middleware
	testUserId := primitive.NewObjectID()
	testUser = &model.User{
		ID:        testUserId,
		UserAlias: testUserId.Hex()[:8],
		Email:     "test@domain.com",
		Password:  "03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4",
	}

	// Prepare authentication header for requests that require it
	testToken := model.Token{}
	err := testToken.CreateToken(testUser)
	if err != nil {
		fmt.Println(err)
		return
	}
	authHeader = map[string]string{
		"Authorization": "Bearer " + testToken.TokenString,
	}

	m.Run()
}

func TestAllEndpointsNormal(t *testing.T) {
	// Start new test server
	mg := NewMockDB()
	apiServer = NewTestServer(mg, middleware.AuthHandler(mg))
	defer apiServer.Close()

	var rr *httptest.ResponseRecorder

	// Create credentials
	credBody, err := json.Marshal(map[string]string{
		"email":    "test",
		"password": "1234",
	})
	if err != nil {
		t.Errorf("Unexpected error")
	}

	// POST /auth/signup should return 201 (Created) and create new user
	rr = sendRequestToMockServer(t, mg, "POST", "/auth/signup", bytes.NewBuffer(credBody), nil)
	assert.Equal(t, http.StatusCreated, rr.Code)

	// POST /auth/token should return 200 (OK) for valid credentials
	rr = sendRequestToMockServer(t, mg, "POST", "/auth/token", bytes.NewBuffer(credBody), nil)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Parse token response
	authResponse := map[string]string{}
	err = json.NewDecoder(rr.Body).Decode(&authResponse)
	if err != nil {
		t.Errorf("Unexpected error")
	}

	// Create authorization header with token
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + authResponse["token"]

	// HEAD /auth/token should return 200 (OK) for valid tokens
	rr = sendRequestToMockServer(t, mg, "HEAD", "/auth/token", nil, headers)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Create a test ride
	rideBody, err := json.Marshal(map[string]interface{}{
		"rideName": "Unit test ride",
		"rideDate": "2023-04-13T11:23:45Z",
		"rideData": []map[string]interface{}{
			{
				"timestamp":      "2023-04-13T11:23:46Z",
				"gyroscopeX":     0.001,
				"gyroscopeY":     -0.002,
				"gyroscopeZ":     0.003,
				"accelerometerX": 0.004,
				"accelerometerY": 0.005,
				"accelerometerZ": -0.006,
				"locationLat":    25.105497,
				"locationLong":   121.597366,
			},
			{
				"timestamp":      "2023-04-13T11:23:50Z",
				"gyroscopeX":     0.002,
				"gyroscopeY":     -0.004,
				"gyroscopeZ":     0.006,
				"accelerometerX": 0.008,
				"accelerometerY": 0.010,
				"accelerometerZ": -0.012,
				"locationLat":    25.105597,
				"locationLong":   121.598366,
			},
		},
	})
	if err != nil {
		t.Errorf("Unexpected error")
	}

	// POST /rides should return 201 (Created) and create new ride
	rr = sendRequestToMockServer(t, mg, "POST", "/rides", bytes.NewBuffer(rideBody), headers)
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Parse ride response
	rideResponse := map[string]interface{}{}
	err = json.NewDecoder(rr.Body).Decode(&rideResponse)
	if err != nil {
		t.Errorf("Unexpected error")
	}

	// POST /rides/{id} should return 200 (OK)
	rideId := rideResponse["_id"].(string)
	rr = sendRequestToMockServer(t, mg, "GET", "/rides/"+rideId, nil, headers)
	assert.Equal(t, http.StatusOK, rr.Code)

	// GET /profile/score should return 200 (OK)
	rr = sendRequestToMockServer(t, mg, "GET", "/profile/score", nil, headers)
	assert.Equal(t, http.StatusOK, rr.Code)

	// GET /profile/stats should return 200 (OK)
	rr = sendRequestToMockServer(t, mg, "GET", "/profile/stats", nil, headers)
	assert.Equal(t, http.StatusOK, rr.Code)

	// DELETE /rides/{id} should return 204 (No Content)
	rr = sendRequestToMockServer(t, mg, "DELETE", "/rides/"+rideId, nil, headers)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestRides(t *testing.T) {
	// Start new test server
	mg := NewMockDB()
	apiServer = NewTestServer(mg, middleware.AuthHandler(mg))
	defer apiServer.Close()

	var rr *httptest.ResponseRecorder

	// GET /rides would return 401 (Unauthorized) with if user does not exist
	rr = sendRequestToMockServer(t, mg, "GET", "/rides", nil, authHeader)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Create test user in database
	err := createTestUser(mg, testUser)
	if err != nil {
		fmt.Println(err)
		return
	}

	// GET /rides would return 200 (OK) with empty profile
	rr = sendRequestToMockServer(t, mg, "GET", "/rides", nil, authHeader)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestProfile(t *testing.T) {
	// Start new test server
	mg := NewMockDB()
	apiServer = NewTestServer(mg, middleware.AuthHandler(mg))
	defer apiServer.Close()

	var rr *httptest.ResponseRecorder

	// GET /profile/score would return 401 (Unauthorized) with if user does not exist
	rr = sendRequestToMockServer(t, mg, "GET", "/profile/score", nil, authHeader)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// GET /profile/stats would return 401 (Unauthorized) with if user does not exist
	rr = sendRequestToMockServer(t, mg, "GET", "/profile/stats", nil, authHeader)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Create test user in database
	err := createTestUser(mg, testUser)
	if err != nil {
		fmt.Println(err)
		return
	}

	// GET /profile/score would return 200 (OK) with empty profile
	rr = sendRequestToMockServer(t, mg, "GET", "/profile/score", nil, authHeader)
	assert.Equal(t, http.StatusOK, rr.Code)

	// GET /profile/stat would return 200 (OK) with empty profile
	rr = sendRequestToMockServer(t, mg, "GET", "/profile/stats", nil, authHeader)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Create test ride in database
	testRide, err := createTestRide(mg, testUser)
	if err != nil {
		fmt.Println(err)
		return
	}

	// GET /profile/score would return 200 (OK)
	rr = sendRequestToMockServer(t, mg, "GET", "/profile/score", nil, authHeader)
	var score model.RideScore
	_ = unmarshalJSON(rr, &score)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, testRide.RideScore.Overall, score.Overall)

	// GET /profile/stat would return 200 (OK)
	rr = sendRequestToMockServer(t, mg, "GET", "/profile/stats", nil, authHeader)
	var stats model.UserStats
	_ = unmarshalJSON(rr, &stats)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, 1, stats.RidesCount)
}
