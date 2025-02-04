package main

import (
	"errors"
	"net/http"

	"github.com/Blue-Davinci/SocialAid/internal/data"
	"github.com/Blue-Davinci/SocialAid/internal/validator"
)

func (app *application) createNewProgramdHandler(w http.ResponseWriter, r *http.Request) {
	// make an input struct that will will hold the inputs we require
	var input struct {
		Name        string `json:"name"`
		Category    string `json:"category"`
		Description string `json:"description"`
	}
	// read the request to the input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		// output a bad request error
		app.badRequestResponse(w, r, err)
		return
	}
	// create a new Program struct and read the input struct to it
	program := &data.Program{
		Name:        input.Name,
		Category:    input.Category,
		Description: input.Description,
	}
	// validate the program struct
	// validate
	v := validator.New()
	if data.ValidateProgram(v, program); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// we are good now, lets create the program
	err = app.models.Program.CreateNewProgram(program)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateProgram):
			v.AddError("name", "a program with this name already exists")
			app.conflictResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// output to client
	err = app.writeJSON(w, http.StatusCreated, envelope{"program": program}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// UpdateProgramByIdHandler() is a handler that updates a program by its ID
// This route supports partial updates. What we do is, we acquire the program's
// ID from the URL as a parameter, check if the program exists, read the request
// using pointer, check which fields have been updated, validate the updated fields,
// and update the program.
func (app *application) updateProgramByIdHandler(w http.ResponseWriter, r *http.Request) {
	// get the program ID from the URL
	id, err := app.readIDParam(r, "programID")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// validate the program ID
	v := validator.New()
	if data.ValidateURLID(v, id, "programID"); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// input struct to hold the user supplied fields
	var input struct {
		Name        *string `json:"name"`
		Category    *string `json:"category"`
		Description *string `json:"description"`
	}
	// read the request to the input struct
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// get the program by its ID
	program, err := app.models.Program.GetProgramById(int32(id))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrProgramDoesNotExist):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// check if the user has updated any field
	if input.Name != nil {
		program.Name = *input.Name
	}
	if input.Category != nil {
		program.Category = *input.Category
	}
	if input.Description != nil {
		program.Description = *input.Description
	}
	// validate the updated fields
	if data.ValidateProgram(v, program); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// update the program
	err = app.models.Program.UpdateProgramById(program)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// output to client
	err = app.writeJSON(w, http.StatusOK, envelope{"program": program}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
