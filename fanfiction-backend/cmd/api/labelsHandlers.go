package main

import (
	"fmt"
	"net/http"

	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/data"
)

// createLabelHandler is the handler used in creating labels
func (app *application) createLabelHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a label")
}

// showLabelHandler is the handler used to show a specific label based on labelID
func (app *application) showLabelHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	label := data.Label{
		ID: id,
		Title: "Student",
		Version: 1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"label": label}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
