package main

import (
	"errors"
	"net/http"

	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/data"
	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/validator"
)

/**
* TODO Create an appropriate error response when creating a label which has given its own id in sublabel/blacklist
* ? When passing a label which doesnt exist into create of label - sublabel/blacklist, it is quietly ignored. Error message?
* ? While creating sublabels/blacklist, should the array have the id's of the labels or the names.
* ? If names, helper function getLabelIDbyName is required
 */

// createLabelHandler is the handler used in creating labels
func (app *application) createLabelHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string  `json:"name"`
		StoryID   int64   `json:"story_id"`
		SubLabels []int64 `json:"sublabels"`
		Blacklist []int64 `json:"blacklist"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	label := &data.Label{
		Name:      input.Name,
		StoryID:   input.StoryID,
		SubLabels: input.SubLabels,
		Blacklist: input.Blacklist,
	}

	v := validator.New()

	if data.ValidateLabel(v, label); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Labels.Create(label)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateLabel):
			v := validator.New()
			v.AddError("label", "a label with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"label": label}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// showLabelHandler is the handler used to show a specific label based on labelID
func (app *application) showLabelHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	label, err := app.models.Labels.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"label": label}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listLabelsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		StoryID int64 `json:"story_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	labels, err := app.models.Labels.GetAllLabels(input.StoryID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"labels": labels}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listLabelsbyCharacterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Character_ID int64 `json:"character_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	labels, err := app.models.Labels.GetAllForCharacter(input.Character_ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"labels": labels}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateLabelHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		Name string `json:"name"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	label, err := app.models.Labels.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	label.Name = input.Name

	err = app.models.Labels.Update(label)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateLabel):
			v := validator.New()
			v.AddError("label", "a label with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"Label": label}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// deleteLabelHandler is the handler used to delete labels based on labelID
func (app *application) deleteLabelHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Labels.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Label successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
