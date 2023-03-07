package main

import (
	"flag"
	"net/http"

	"github.com/ktalexcheng/trailbrake_api/api/router"
	"github.com/ktalexcheng/trailbrake_api/config"
	"github.com/ktalexcheng/trailbrake_api/util"
	"github.com/rs/zerolog/log"
)

func main() {
	debug := flag.Bool("debug", false, "Running in debug mode")
	flag.Parse()

	// Initialize global logger settings
	util.InitLogger(*debug)

	// Initialize new database client
	mongoClient, err := util.NewMongoClient(config.DbConnection, config.DbName, config.DbCertificate)
	if err != nil {
		panic("Unable to connect to MongoDB.")
	}

	// Get router and start service
	r := router.Router(mongoClient)
	log.Info().Msg("Started Trailbrake API service")
	http.ListenAndServe(":8080", r)
}
