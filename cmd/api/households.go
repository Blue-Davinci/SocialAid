package main

import (
	"errors"
	"net/http"

	"github.com/Blue-Davinci/SocialAid/internal/data"
	"github.com/Blue-Davinci/SocialAid/internal/validator"
)

// createNewHouseHoldHandler() is a handler that creates a new house hold
// We read the request, validate the input, if successful we create a new house hold
// and output the house hold to the client
func (app *application) createNewHouseHoldHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ProgramID     int32  `json:"program_id"`
		GeoLocationID int32  `json:"geo_location_id"`
		Name          string `json:"name"`
	}
	// read the request to the input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// create a new HouseHold struct and read the input struct to it
	houseHold := &data.HouseHold{
		ProgramID:     input.ProgramID,
		GeoLocationID: input.GeoLocationID,
		Name:          input.Name,
	}
	// validate the house hold struct
	v := validator.New()
	if data.ValidateHouseHold(v, houseHold); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// we are good now, lets create the house hold
	err = app.models.HouseHold.CreateNewHouseHold(houseHold)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrGeoLocationDoesNotExist):
			v.AddError("geo_location_id", "geo location does not exist")
			app.conflictResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrProgramDoesNotExist):
			v.AddError("program_id", "program does not exist")
			app.conflictResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// output to client
	err = app.writeJSON(w, http.StatusCreated, envelope{"house_hold": houseHold}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
