package main

import (
	"errors"
	"net/http"

	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/data"
	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/validator"
)

func (app *application) createTimelineHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Story_ID int64 `json:"story_id"`
		Name string `json:"name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	
	timeline := &data.Timeline{
		Story_ID: input.Story_ID,
		Name: input.Name,	
	}

	v := validator.New()
	if data.ValidateTimeline(v, timeline); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Timelines.Insert(timeline)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateTimeline):
			v.AddError("name", "a timeline with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"timeline": timeline}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
