package router

import (
	"github.com/go-chi/chi"

	"github.com/ktalexcheng/trailbrake_api/api/handler"
	"github.com/ktalexcheng/trailbrake_api/api/middleware"
	"github.com/ktalexcheng/trailbrake_api/util"
)

// Initializes chi.NewRouter() and map handler to endpoints
func Router(mg *util.MongoClient) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.LogHandler)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/token", handler.GetToken(mg))
		r.Post("/signup", handler.CreateNewUser(mg))
	})

	r.Route("/rides", func(r chi.Router) {
		r.Use(middleware.AuthHandler(mg))
		r.Post("/", handler.SaveRideData(mg))
		r.Get("/", handler.GetAllRideRecords(mg))
		r.Get("/{rideId}", handler.GetRideData(mg))
		r.Delete("/{rideId}", handler.DeleteRideData(mg))
	})

	r.Route("/profile", func(r chi.Router) {
		r.Use(middleware.AuthHandler(mg))
		r.Get("/score", handler.GetUserScore(mg))
		r.Get("/stats", handler.GetUserStats(mg))
	})

	return r
}
