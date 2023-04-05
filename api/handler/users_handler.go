package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/ktalexcheng/trailbrake_api/api/middleware"
	"github.com/ktalexcheng/trailbrake_api/api/model"
	"github.com/ktalexcheng/trailbrake_api/util"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserScore(mg *util.MongoClient) http.HandlerFunc {
	return middleware.ErrorHandler(
		func(w http.ResponseWriter, r *http.Request) error {
			param := chi.URLParam(r, "useRecentRides")

			var useRecentRides int
			var err error
			if param != "" {
				useRecentRides, err = strconv.Atoi(param)
				if err != nil {
					return err
				}

				if useRecentRides < 0 {
					return errors.New("useRecentRides must be a non-negative integer")
				}
			} else {
				useRecentRides = 10
			}

			userScore, err := doGetUserScore(r.Context(), mg, useRecentRides)
			if err != nil {
				return err
			}

			if userScore == nil {
				err = util.HTTPWriteJSONResponse(w, http.StatusNotFound, &util.JSONResponse{
					Message: "no rides found for user",
				})
				if err != nil {
					return err
				}
				return nil
			}

			userScoreJson, err := json.Marshal(userScore)
			if err != nil {
				return err
			}

			_, err = w.Write(userScoreJson)
			if err != nil {
				return err
			}

			return nil
		},
	)
}

func doGetUserScore(ctx context.Context, mg *util.MongoClient, useRecentRides int) (*model.RideScore, error) {
	token := ctx.Value(model.TokenKey).(model.Token)
	userId, err := primitive.ObjectIDFromHex(token.Subject)
	if err != nil {
		return nil, err
	}

	pipeline := make([]bson.M, 0)

	if useRecentRides > 0 {
		pipeline = append(pipeline, bson.M{
			"$limit": useRecentRides,
		})
	}

	pipeline = append(pipeline, []bson.M{
		{
			"$match": bson.M{
				"userId": userId,
			},
		},
		{
			"$sort": bson.M{
				"rideDate": -1,
			},
		},
		{
			"$group": bson.M{
				"_id":            nil,
				"_totalDuration": bson.M{"$sum": "$rideMeta.duration"},
				"_sumProdOverall": bson.M{
					"$sum": bson.M{
						"$multiply": []string{
							"$rideScore.overall",
							"$rideMeta.duration",
						},
					},
				},
				"_sumProdAcceleration": bson.M{
					"$sum": bson.M{
						"$multiply": []string{
							"$rideScore.acceleration",
							"$rideMeta.duration",
						},
					},
				},
				"_sumProdBraking": bson.M{
					"$sum": bson.M{
						"$multiply": []string{
							"$rideScore.braking",
							"$rideMeta.duration",
						},
					},
				},
				"_sumProdCornering": bson.M{
					"$sum": bson.M{
						"$multiply": []string{
							"$rideScore.cornering",
							"$rideMeta.duration",
						},
					},
				},
				"_sumProdSpeed": bson.M{
					"$sum": bson.M{
						"$multiply": []string{
							"$rideScore.speed",
							"$rideMeta.duration",
						},
					},
				},
			},
		},
		{
			"$project": bson.M{
				"overall": bson.M{
					"$divide": []string{"$_sumProdOverall", "$_totalDuration"},
				},
				"acceleration": bson.M{
					"$divide": []string{"$_sumProdAcceleration", "$_totalDuration"},
				},
				"braking": bson.M{
					"$divide": []string{"$_sumProdBraking", "$_totalDuration"},
				},
				"cornering": bson.M{
					"$divide": []string{"$_sumProdCornering", "$_totalDuration"},
				},
				"speed": bson.M{
					"$divide": []string{"$_sumProdSpeed", "$_totalDuration"},
				},
			},
		},
	}...)

	// rideRecordsCol := mg.MongoDB.Collection("rideRecords")
	rideRecordsCol := mg.RideRecordsColl
	cur, err := rideRecordsCol.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	var userScore []model.RideScore
	err = cur.All(context.TODO(), &userScore)
	if err != nil {
		return nil, err
	}

	if len(userScore) == 0 {
		return nil, nil
	}

	return &userScore[0], nil
}

func GetUserStats(mg *util.MongoClient) http.HandlerFunc {
	return middleware.ErrorHandler(
		func(w http.ResponseWriter, r *http.Request) error {
			userStats, err := doGetUserStats(r.Context(), mg)
			if err != nil {
				return err
			}

			if userStats == nil {
				err = util.HTTPWriteJSONResponse(w, http.StatusNotFound, &util.JSONResponse{
					Message: "no rides found for user",
				})
				if err != nil {
					return err
				}
				return nil
			}

			userStatsJson, err := json.Marshal(userStats)
			if err != nil {
				return err
			}

			_, err = w.Write(userStatsJson)
			if err != nil {
				return err
			}

			return nil
		},
	)
}

func doGetUserStats(ctx context.Context, mg *util.MongoClient) (*model.UserStats, error) {
	token := ctx.Value(model.TokenKey).(model.Token)
	userId, err := primitive.ObjectIDFromHex(token.Subject)
	if err != nil {
		return nil, err
	}

	pipeline := make([]bson.M, 0)

	pipeline = append(pipeline, []bson.M{
		{
			"$match": bson.M{
				"userId": userId,
			},
		},
		{"$group": bson.M{
			"_id":             nil,
			"ridesCount":      bson.M{"$count": bson.M{}},
			"totalDistance":   bson.M{"$sum": "$rideMeta.distance"},
			"totalRideTime":   bson.M{"$sum": "$rideMeta.duration"},
			"maxAcceleration": bson.M{"$max": "$rideMeta.maxAcceleration"},
		},
		},
	}...)

	// rideRecordsCol := mg.MongoDB.Collection("rideRecords")
	rideRecordsCol := mg.RideRecordsColl
	cur, err := rideRecordsCol.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}

	var userStats []model.UserStats
	err = cur.All(context.TODO(), &userStats)
	if err != nil {
		return nil, err
	}

	if len(userStats) == 0 {
		return nil, nil
	}

	return &userStats[0], nil
}

func CreateNewUser(mg *util.MongoClient) http.HandlerFunc {
	return middleware.ErrorHandler(
		func(w http.ResponseWriter, r *http.Request) error {
			log.Info().Msg("Creating new user")

			var user model.User
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				return err
			}

			if user.Email == "" || user.Password == "" {
				err = util.HTTPWriteJSONResponse(w, http.StatusBadRequest, &util.JSONResponse{
					Message: "'email' and 'password' must not be blank.",
				})
				if err != nil {
					return err
				}
			}

			token, err := doCreateNewUser(mg, &user)
			if err != nil {
				return err
			}

			// Return the token in the response
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(map[string]string{
				"token": token.TokenString,
			})
			if err != nil {
				return err
			}

			log.Info().Msg("Create new user successful")
			return nil
		},
	)
}

func doCreateNewUser(mg *util.MongoClient, user *model.User) (*model.Token, error) {
	// usersCol := mg.MongoDB.Collection("users")
	usersCol := mg.UsersColl

	count, err := usersCol.CountDocuments(context.TODO(), bson.D{{Key: "email", Value: (*user).Email}})
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("user already exists")
	}

	user.ID = primitive.NewObjectID()
	user.UserAlias = user.ID.Hex()[:8]
	_, err = usersCol.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}

	var token model.Token
	err = token.CreateToken(user)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func GetToken(mg *util.MongoClient) http.HandlerFunc {
	return middleware.ErrorHandler(
		func(w http.ResponseWriter, r *http.Request) error {
			if r.ContentLength == 0 {
				err := util.HTTPWriteJSONResponse(w, http.StatusUnauthorized, &util.JSONResponse{
					Message: "missing credentials",
				})
				if err != nil {
					return err
				}

				return nil
			}

			var user model.User
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				return err
			}

			token, err := doGetToken(mg, &user)
			if err != nil {
				return err
			}
			// Invalid credentials
			if token == nil {
				err = util.HTTPWriteJSONResponse(w, http.StatusUnauthorized, &util.JSONResponse{
					Message: "invalid credentials",
				})
				if err != nil {
					return err
				}

				return nil
			}

			response, err := json.Marshal(map[string]interface{}{
				"token": token.TokenString,
			})
			if err != nil {
				return err
			}

			w.WriteHeader(http.StatusOK)
			_, err = w.Write(response)
			if err != nil {
				return err
			}

			return nil
		},
	)
}

func doGetToken(mg *util.MongoClient, user *model.User) (*model.Token, error) {
	userExists, err := user.CheckUserExists(mg)
	if err != nil {
		return nil, err
	}

	if userExists {
		userValid, err := user.ValidateUserPass(mg)
		if err != nil {
			return nil, err
		}

		if userValid {
			var token model.Token

			err = token.CreateToken(user)
			if err != nil {
				return nil, err
			}

			return &token, nil
		}
	}

	return nil, nil
}
