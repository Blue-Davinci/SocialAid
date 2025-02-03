package main

import (
	"net/http"

	"go.uber.org/zap"
)

func (app *application) createNewProgramdHandler(w http.ResponseWriter, r *http.Request) {
	// create a new program
	// output to client
	app.logger.Info("createNewProgramdHandler", zap.String("Task", "Create a new program"))
}
