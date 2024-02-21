package main

import (
	"errors"
	"net/http"

	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/data"
	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/validator"
	"github.com/google/uuid"
)

func (app *application) createEventHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Character_ID string `json:"character_id"`
		// EventTime    string `json:"event_time"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Details     string `json:"details"`
		Index       int    `json:"index"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// eventTime, err := time.Parse("2006-01-02 15:04", input.EventTime)
	// if err != nil {
	// 	app.serverErrorResponse(w, r, err)
	// 	return
	// }

	char_uuid, err := uuid.Parse(input.Character_ID)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	event := &data.Event{
		Character_ID: char_uuid,
		// EventTime:    eventTime,
		Title:       input.Title,
		Description: input.Description,
		Details:     input.Details,
		Index:       input.Index,
	}

	v := validator.New()
	if data.ValidateEvent(v, event); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Events.Insert(event)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEvent):
			v.AddError("name", "a event with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"event": event}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"event": event}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listEventHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Character_ID string `json:"character_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	char_uuid, err := uuid.Parse(input.Character_ID)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	events, err := app.models.Events.GetForCharacter(char_uuid)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"events": events}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listStoryEventHandler(w http.ResponseWriter, r *http.Request) {
	// var input struct {
	// 	Story_ID int64 `json:"story_id"`
	// }

	// err := app.readJSON(w, r, &input)
	// if err != nil {
	// 	app.badRequestResponse(w, r, err)
	// 	return
	// }

	// events, err := app.models.Events.GetForStory(input.Story_ID)

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	
	events, err := app.models.Events.GetForStory(id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// err = app.writeJSON(w, http.StatusOK, envelope{"story_id": input.Story_ID, "story_events": events}, nil)
	err = app.writeJSON(w, http.StatusOK, envelope{"story_id": id, "story_events": events}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		Character_ID *string `json:"character_id"`
		// EventTime    *string `json:"event_time"`
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Details     *string `json:"details"`
		Index       *int    `json:"index"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	char_uuid, err := uuid.Parse(*input.Character_ID)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	event, err := app.models.Events.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	oldIndex := event.Index
	oldCharId := event.Character_ID

	if input.Character_ID != nil {
		event.Character_ID = char_uuid
	}

	// if input.EventTime != nil {
	// 	event.EventTime, err = time.Parse("2006-01-02 15:04", *input.EventTime)
	// 	if err != nil {
	// 		app.serverErrorResponse(w, r, err)
	// 		return
	// 	}
	// }
	if input.Title != nil {
		event.Title = *input.Title
	}
	if input.Description != nil {
		event.Description = *input.Description
	}
	if input.Details != nil {
		event.Details = *input.Details
	}
	if input.Index != nil {
		event.Index = *input.Index
	}

	v := validator.New()
	if data.ValidateEvent(v, event); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Events.Update(event, oldIndex, oldCharId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEvent):
			v := validator.New()
			v.AddError("events", "a event with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"Event": event}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Character_ID string `json:"character_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	uuid, err := uuid.Parse(input.Character_ID)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Events.Delete(id, uuid)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Event successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
