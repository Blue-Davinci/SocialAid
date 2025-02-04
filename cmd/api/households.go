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

// getHouseHoldInformationHandler() is a handler that gets the information of a house hold
// We get the ID from the request as a URL parameter, validate the input. If everything is okay,
// we pass down the house hold information to the client
func (app *application) getHouseHoldInformationHandler(w http.ResponseWriter, r *http.Request) {
	// get the house hold id from the URL parameter
	houseHoldID, err := app.readIDParam(r, "householdID")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// validate the household ID
	v := validator.New()
	if data.ValidateURLID(v, houseHoldID, "householdID"); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// get the house hold information
	houseHold, err := app.models.HouseHold.GetHouseHoldInformation(int32(houseHoldID), app.config.encryption.key)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrHouseHoldDoesNotExist):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// output to client
	err = app.writeJSON(w, http.StatusOK, envelope{"house_hold": houseHold}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// createNewHouseholdHeadHandler() is a handler that creates a new house hold head
// We read the request, validate the input. If everything is okay, we pass down the
// new household head as well as the encryption key to be saved in the database
func (app *application) createNewHouseholdHeadHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		HouseHoldID int32  `json:"house_hold_id"`
		Name        string `json:"name"`
		NationalID  string `json:"national_id"`
		PhoneNumber string `json:"phone_number"`
		Age         int32  `json:"age"`
	}
	// read the request to the input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// create a new HouseHoldHead struct and read the input struct to it
	houseHoldHead := &data.HouseHoldHead{
		HouseHoldID: input.HouseHoldID,
		Name:        input.Name,
		NationalID:  input.NationalID,
		PhoneNumber: input.PhoneNumber,
		Age:         input.Age,
	}
	// validate the house hold head struct
	v := validator.New()
	if data.ValidateHouseHoldHead(v, houseHoldHead); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// we are good now, lets create the house hold head
	err = app.models.HouseHold.CreateNewHouseholdHead(houseHoldHead, app.config.encryption.key)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrHouseHoldDoesNotExist):
			v.AddError("house_hold_id", "house hold does not exist")
			app.conflictResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrHouseHoldAlreadyExists):
			v.AddError("house_hold_id", "house hold already exists and has a head")
			app.conflictResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// output to client
	err = app.writeJSON(w, http.StatusCreated, envelope{"house_hold_head": houseHoldHead}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// createNewHouseholdMemberHandler() is a handler that creates a new house hold member
// We read the request, validate the input. If everything is okay, we pass down the
// new household member to be saved in the database
func (app *application) createNewHouseholdMemberHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		HouseHoldID int32  `json:"house_hold_id"`
		Name        string `json:"name"`
		Age         int32  `json:"age"`
		Relation    string `json:"relation"`
	}
	// read the request to the input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// create a new HouseHoldMember struct and read the input struct to it
	houseHoldMember := &data.HouseHoldMember{
		HouseHoldID: input.HouseHoldID,
		Name:        input.Name,
		Age:         input.Age,
		Relation:    input.Relation,
	}
	// validate the house hold member struct
	v := validator.New()
	if data.ValidateHouseHoldMember(v, houseHoldMember); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// check if the household head has been created for this household
	// if not, we return an error stating that the household head has not been created yet.
	// This enhances data integrity
	_, err = app.models.HouseHold.GetHouseholdHeadByHouseholdId(houseHoldMember.HouseHoldID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrHouseHoldHeadDoesNotExist):
			v.AddError("house_hold_id", "house hold or its head does not exist")
			app.conflictResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// we are good now, lets create the house hold member
	err = app.models.HouseHold.CreateNewHouseholdMember(houseHoldMember)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrHouseHoldDoesNotExist):
			v.AddError("house_hold_id", "house hold does not exist")
			app.conflictResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrHouseHoldMemberExists):
			v.AddError("house_hold_id", "house hold member already exists")
			app.conflictResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// output to client
	err = app.writeJSON(w, http.StatusCreated, envelope{"house_hold_member": houseHoldMember}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
