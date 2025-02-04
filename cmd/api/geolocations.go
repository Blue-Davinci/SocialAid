package main

import (
	"errors"
	"net/http"

	"github.com/Blue-Davinci/SocialAid/internal/data"
	"github.com/Blue-Davinci/SocialAid/internal/validator"
)

// createNewGeoLocationHandler() is a handler that creates a new geo location
// We read the request, validate the input, if successful we create a new geo location
// and output the geo location to the client
func (app *application) createNewGeoLocationHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		County      string `json:"county"`
		SubCounty   string `json:"sub_county"`
		Location    string `json:"location"`
		SubLocation string `json:"sub_location"`
	}
	// read the request to the input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// create a new GeoLocation struct and read the input struct to it
	geoLocation := &data.GeoLocation{
		County:      input.County,
		SubCounty:   input.SubCounty,
		Location:    input.Location,
		SubLocation: input.SubLocation,
	}
	// validate the geo location struct
	v := validator.New()
	if data.ValidateGeoLocation(v, geoLocation); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// we are good now, lets create the geo location
	err = app.models.GeoLocation.CreateNewGeoLocation(geoLocation)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateGeoLocation):
			v.AddError("county", "a geo location with this sub-location already exists")
			app.conflictResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// output to client
	err = app.writeJSON(w, http.StatusCreated, envelope{"geo_location": geoLocation}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
