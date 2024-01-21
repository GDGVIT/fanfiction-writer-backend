package main

import (
	"errors"
	"net/http"

	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/data"
	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/validator"
)

func (app *application) createTimelineHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Story_ID int64  `json:"story_id"`
		Name     string `json:"name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	timeline := &data.Timeline{
		Story_ID: input.Story_ID,
		Name:     input.Name,
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
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
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

func (app *application) getTimelineHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Story_ID int64 `json:"story_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	timeline, err := app.models.Timelines.Get(input.Story_ID, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"timeline": timeline}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listTimelineHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Story_ID int64 `json:"story_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	timelines, err := app.models.Timelines.GetForStory(input.Story_ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"timelines": timelines}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateTimelineHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		Story_ID int64  `json:"story_id"`
		Name     string `json:"name"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	timeline, err := app.models.Timelines.Get(input.Story_ID, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	timeline.Name = input.Name

	err = app.models.Timelines.Update(timeline)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateTimeline):
			v := validator.New()
			v.AddError("timeline", "a timeline with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"Timeline": timeline}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) deleteTimelineHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Story_ID int64 `json:"story_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Timelines.Delete(input.Story_ID, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Timeline successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
