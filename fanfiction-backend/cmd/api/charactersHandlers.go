package main

import (
	"errors"
	"net/http"

	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/data"
	"github.com/GDGVIT/fanfiction-writer-backend/fanfiction-backend/internal/validator"
	"github.com/google/uuid"
)

func (app *application) createCharacterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Story_ID    int64  `json:"story_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Index       int    `json:"index"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	character := &data.Character{
		Story_ID:    input.Story_ID,
		Description: input.Description,
		Name:        input.Name,
		Index:       input.Index,
	}

	v := validator.New()
	if data.ValidateCharacter(v, character); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Characters.Insert(character)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateCharacter):
			v.AddError("name", "a character with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"character": character}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) createCharLabelHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Character_ID string  `json:"character_id"`
		Label_ID     []int64 `json:"label_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	id, err := uuid.Parse(input.Character_ID)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Characters.InsertCharLabels(id, input.Label_ID...)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrDuplicateLabel):
			v := validator.New()
			v.AddError("charlabel", "a character-label entry with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"character": "character label(s) added."}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getCharacterHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	character, err := app.models.Characters.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"character": character}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listCharacterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Story_ID int64 `json:"story_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	characters, err := app.models.Characters.GetForStory(input.Story_ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"characters": characters}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) listCharacterByLabelsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Label_ID int64 `json:"label_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	characters, err := app.models.Characters.GetAllForLabel(input.Label_ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"characters": characters}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateCharacterHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		Story_ID    *int64  `json:"story_id"`
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Index       *int    `json:"index"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	character, err := app.models.Characters.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	oldIndex := character.Index

	if input.Name != nil {
		character.Name = *input.Name
	}
	if input.Description != nil {
		character.Description = *input.Description
	}
	if input.Index != nil {
		character.Index = *input.Index
	}

	v := validator.New()
	if data.ValidateCharacter(v, character); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Characters.Update(character, oldIndex)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateCharacter):
			v := validator.New()
			v.AddError("character", "a character with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"character": character}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) deleteCharacterHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam(r)
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

	err = app.models.Characters.Delete(input.Story_ID, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Character successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteCharLabelHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Character_ID string  `json:"character_id"`
		Label_ID     []int64 `json:"label_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	id, err := uuid.Parse(input.Character_ID)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Characters.DeleteCharLabels(id, input.Label_ID...)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"character": "character label(s) deleted."}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
