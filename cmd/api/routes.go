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

	// Mount to our Versioning router
	router.Mount("/v1", v1Router)
	return router
}

func (app *application) programRoutes() http.Handler {
	router := chi.NewRouter()
	router.Post("/", app.createNewProgramdHandler)
	return router
}
