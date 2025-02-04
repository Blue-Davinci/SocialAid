package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   app.config.cors.trustedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// create main router
	v1Router := chi.NewRouter()
	v1Router.Mount("/programs", app.programRoutes())
	v1Router.Mount("/geo_locations", app.geoLocationRoutes())
	v1Router.Mount("/house_holds", app.houseHoldRoutes())

	// Mount to our Versioning router
	router.Mount("/v1", v1Router)
	return router
}

// programRoutes() is a route handler responsible for all program routes
func (app *application) programRoutes() http.Handler {
	router := chi.NewRouter()
	router.Post("/", app.createNewProgramdHandler)
	router.Patch("/{programID}", app.updateProgramByIdHandler)
	return router
}

// geoLocationRoutes() is a route handler responsible for all geo location routes
func (app *application) geoLocationRoutes() http.Handler {
	router := chi.NewRouter()
	router.Post("/", app.createNewGeoLocationHandler)
	return router
}

// houseHoldRoutes() is a route handler responsible for all house hold routes
func (app *application) houseHoldRoutes() http.Handler {
	router := chi.NewRouter()
	router.Get("/{householdID}", app.getHouseHoldInformationHandler) // GET request with a parameter
	router.Post("/", app.createNewHouseHoldHandler)

	// household head
	router.Post("/head", app.createNewHouseholdHeadHandler)
	// household member
	router.Post("/member", app.createNewHouseholdMemberHandler)
	return router
}
