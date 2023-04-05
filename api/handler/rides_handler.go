package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/ktalexcheng/trailbrake_api/api/middleware"
	"github.com/ktalexcheng/trailbrake_api/api/model"
	"github.com/ktalexcheng/trailbrake_api/util"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteRideData(mg *util.MongoClient) http.HandlerFunc {
	return middleware.ErrorHandler(
		func(w http.ResponseWriter, r *http.Request) error {
			log.Info().Msg("Deleting ride data")

			rideId := chi.URLParam(r, "rideId")

			err := doDeleteRideData(mg, rideId)
			if err != nil {
				return err
			}

			log.Info().Msg("Delete ride data successful")
			return nil
		},
	)
}

func doDeleteRideData(mg *util.MongoClient, rideId string) error {
	// rideDataCol := mg.MongoDB.Collection("rideData")
	// rideRecordsCol := mg.MongoDB.Collection("rideRecords")
	rideDataCol := mg.RideDataColl
	rideRecordsCol := mg.RideRecordsColl

	objId, err := primitive.ObjectIDFromHex(rideId)
	if err != nil {
		return err
	}

	rideDataFilter := bson.D{
		{Key: "rideRecordId", Value: objId},
	}
	rideRecFilter := bson.D{
		{Key: "_id", Value: objId},
	}

	// TODO: Check userId matches before delete

	_, err = rideDataCol.DeleteMany(context.TODO(), rideDataFilter)
	if err != nil {
		return err
	}

	_, err = rideRecordsCol.DeleteOne(context.TODO(), rideRecFilter)
	if err != nil {
		return err
	}

	return nil
}

func SaveRideData(mg *util.MongoClient) http.HandlerFunc {
	return middleware.ErrorHandler(
		func(w http.ResponseWriter, r *http.Request) error {
			log.Info().Msg("Saving ride data")

			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			if err != nil {
				return err
			}

			rideRecord, err := doSaveRideData(r.Context(), mg, &body)
			if err != nil {
				return err
			}

			rideRecordJson, err := json.Marshal(rideRecord)
			if err != nil {
				return err
			}
			_, err = w.Write(rideRecordJson)
			if err != nil {
				return err
			}

			log.Info().Msg("Save ride data successful")
			return nil
		},
	)
}

func doSaveRideData(ctx context.Context, mg *util.MongoClient, body *map[string]interface{}) (*model.RideRecord, error) {
	newRideRecordID := primitive.NewObjectID()
	tsStr := (*body)["rideDate"].(string)
	timestamp, err := time.Parse(time.RFC3339, tsStr)
	if err != nil {
		return nil, err
	}

	var rideData []interface{}
	for _, d := range (*body)["rideData"].([]interface{}) {
		datum := d.(map[string]interface{})
		tsStr = datum["timestamp"].(string)
		timestamp, err = time.Parse(time.RFC3339, tsStr)
		if err != nil {
			return nil, err
		}

		rideData = append(rideData, model.RideDatum{
			Timestamp:      timestamp,
			RideRecordID:   newRideRecordID,
			LocationLat:    datum["locationLat"].(float64),
			LocationLong:   datum["locationLong"].(float64),
			AccelerometerX: datum["accelerometerX"].(float64),
			AccelerometerY: datum["accelerometerY"].(float64),
			AccelerometerZ: datum["accelerometerZ"].(float64),
			GyroscopeX:     datum["gyroscopeX"].(float64),
			GyroscopeY:     datum["gyroscopeY"].(float64),
			GyroscopeZ:     datum["gyroscopeZ"].(float64),
		})
	}

	judgeScore, err := doJudgeRideScore(&rideData)
	if err != nil {
		return nil, err
	}

	token := ctx.Value(model.TokenKey).(model.Token)
	userId, err := primitive.ObjectIDFromHex(token.Subject)
	if err != nil {
		return nil, err
	}
	rideMeta := (*judgeScore).RideMeta

	newRideRecord := model.RideRecord{
		ID:        newRideRecordID,
		RideName:  (*body)["rideName"].(string),
		RideDate:  timestamp,
		RideScore: (*judgeScore).RideScore,
		RideMeta:  rideMeta,
		UserID:    userId,
	}

	// rideRecordsCol := mg.MongoDB.Collection("rideRecords")
	// rideDataCol := mg.MongoDB.Collection("rideData")
	rideRecordsCol := mg.RideRecordsColl
	rideDataCol := mg.RideDataColl

	_, err = rideRecordsCol.InsertOne(context.TODO(), newRideRecord)
	if err != nil {
		return nil, err
	}
	_, err = rideDataCol.InsertMany(context.TODO(), rideData)
	if err != nil {
		return nil, err
	}

	return &newRideRecord, nil
}

func doJudgeRideScore(rideData *[]interface{}) (*model.JudgeResult, error) {
	judgeApiScore := fmt.Sprintf("%s/rideScore", os.Getenv("JUDGE_URL"))
	postBody, _ := json.Marshal(map[string]*[]interface{}{
		"rideData": rideData,
	})
	postBodyBuffer := bytes.NewBuffer(postBody)

	resp, err := http.Post(judgeApiScore, "application/json", postBodyBuffer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if successful
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("judge service error")
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var judgeResult model.JudgeResult
	err = json.Unmarshal(respBody, &judgeResult)
	if err != nil {
		return nil, err
	}

	return &judgeResult, err
}

func GetRideData(mg *util.MongoClient) http.HandlerFunc {
	return middleware.ErrorHandler(
		func(w http.ResponseWriter, r *http.Request) error {
			log.Info().Msg("Getting ride data")

			rideId := chi.URLParam(r, "rideId")

			docs, err := doGetRideData(mg, rideId)
			if err != nil {
				return err
			}

			response, err := json.Marshal(docs)
			if err != nil {
				return err
			}

			_, err = w.Write(response)
			if err != nil {
				return err
			}

			log.Info().Msg("Get ride data successful")
			return nil
		},
	)
}

func doGetRideData(mg *util.MongoClient, rideId string) (*model.Ride, error) {
	objId, err := primitive.ObjectIDFromHex(rideId)

	// rideRecordsCol := mg.MongoDB.Collection("rideRecords")
	// rideDataCol := mg.MongoDB.Collection("rideData")
	rideRecordsCol := mg.RideRecordsColl
	rideDataCol := mg.RideDataColl

	// Fetch ride record
	if err != nil {
		return nil, err
	}
	rideRecFilter := bson.D{
		{Key: "_id", Value: objId},
	}

	// TODO: Check userId matches before get

	var rideRec model.RideRecord
	err = rideRecordsCol.FindOne(context.TODO(), rideRecFilter).Decode(&rideRec)
	if err != nil {
		return nil, err
	}

	// Fetch ride data
	rideDataFilter := bson.D{
		{Key: "metadata.rideRecordID", Value: objId},
	}
	rideDataCur, err := rideDataCol.Find(context.TODO(), rideDataFilter)
	if err != nil {
		return nil, err
	}

	var rideData []model.RideDatum
	err = rideDataCur.All(context.TODO(), &rideData)
	if err != nil {
		return nil, err
	}

	ride := model.Ride{
		ID:        rideRec.ID,
		RideName:  rideRec.RideName,
		RideDate:  rideRec.RideDate,
		RideMeta:  rideRec.RideMeta,
		RideScore: rideRec.RideScore,
		RideData:  rideData,
	}

	return &ride, nil
}

func GetAllRideRecords(mg *util.MongoClient) http.HandlerFunc {
	return middleware.ErrorHandler(
		func(w http.ResponseWriter, r *http.Request) error {
			log.Info().Msg("Getting ride records only")

			docs, err := doGetAllRideRecords(r.Context(), mg)
			if err != nil {
				return err
			}

			response, err := json.Marshal(docs)
			if err != nil {
				return err
			}

			_, err = w.Write(response)
			if err != nil {
				return err
			}

			log.Info().Msg("Get ride records only successful")
			return nil
		},
	)
}

func doGetAllRideRecords(ctx context.Context, mg *util.MongoClient) (*[]model.RideRecord, error) {
	// rideRecordsCol := mg.MongoDB.Collection("rideRecords")
	rideRecordsCol := mg.RideRecordsColl

	token := ctx.Value(model.TokenKey).(model.Token)
	userId, err := primitive.ObjectIDFromHex(token.Subject)
	if err != nil {
		return nil, err
	}

	// pipeline := make([]bson.M, 0)
	// pipeline = append(pipeline, []bson.M{
	// 	{
	// 		"$match": bson.M{
	// 			"userId": userId,
	// 		},
	// 	},
	// }...)

	// cursor, err := rideRecordsCol.Aggregate(context.TODO(), pipeline)
	// if err != nil {
	// 	return nil, err
	// }

	// Fetch ride data for the user only
	rideRecordsFilter := bson.M{
		"userId": userId,
	}
	cursor, err := rideRecordsCol.Find(context.TODO(), rideRecordsFilter)
	if err != nil {
		return nil, err
	}

	var results []model.RideRecord
	err = cursor.All(context.TODO(), &results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}
