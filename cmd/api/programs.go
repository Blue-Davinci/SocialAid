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
