package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/ktalexcheng/trailbrake_api/api/middleware"
	"github.com/ktalexcheng/trailbrake_api/api/router"
	"github.com/ktalexcheng/trailbrake_api/util"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func main() {
	logErr := func(e error) { log.Error().Stack().Err(errors.Wrap(e, "error")).Msg("") }

	debug := flag.Bool("debug", false, "Running in debug mode")
	flag.Parse()
	util.SetupDefaultEnv(*debug)

	// Initialize global logger settings
	util.InitLogger(*debug)

	// Initialize new database client
	var dbCertPath = os.Getenv("MONGO_DB_CERT")
	var dbConnection = os.Getenv("MONGO_DB_CONN")
	var dbName = os.Getenv("MONGO_DB_NAME")
	mongoClient, err := util.NewMongoClient(dbConnection, dbName, dbCertPath)
	if err != nil {
		logErr(err)
		panic("unable to connect to MongoDB.")
	}

	// Get router and start service
	var port = os.Getenv("PORT")
	r := router.Router(mongoClient, middleware.AuthHandler(mongoClient))
	log.Info().Msg("Starting Trailbrake API service")
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		logErr(err)
		panic("unable to start service")
	}
}
