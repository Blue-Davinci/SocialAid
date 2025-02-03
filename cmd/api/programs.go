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
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
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
		app.errorResponse(w, r, http.StatusUnprocessableEntity, v.Errors)
		return
	}
	// we are good now, lets create the program
	err = app.models.Program.CreateNewProgram(program)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateProgram):
			app.errorResponse(w, r, http.StatusConflict, "program name already exists")
		default:
			app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	// output to client
	err = app.writeJSON(w, http.StatusCreated, envelope{"program": program}, nil)
	if err != nil {
		app.errorResponse(w, r, http.StatusInternalServerError, err.Error())
	}
}
