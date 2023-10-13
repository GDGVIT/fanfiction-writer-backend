package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/data"
	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/validator"
)

// createLabelHandler is the handler used in creating labels
func (app *application) createLabelHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string  `json:"name"`
		SubLabels []int64 `json:"sub_labels"`
		Blacklist []int64 `json:"blacklist"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
	}

	label := &data.Label{
		Name: input.Name,
		SubLabels: input.SubLabels,
		Blacklist: input.Blacklist,
	}

	v := validator.New()

	if data.ValidateLabel(v, label);!v.Valid(){
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

// showLabelHandler is the handler used to show a specific label based on labelID
func (app *application) showLabelHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	label := data.Label{
		ID:      id,
		CreatedAt: time.Now(),
		Name:    "Student",
		Version: 1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"label": label}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
