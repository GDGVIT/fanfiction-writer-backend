package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/labels/:id", app.requireActivatedUser(app.showLabelHandler))
	router.HandlerFunc(http.MethodGet, "/v1/labels", app.requireActivatedUser(app.listLabelsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/labels", app.requireActivatedUser(app.createLabelHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/labels/:id", app.requireActivatedUser(app.deleteLabelHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/labels/:id", app.requireActivatedUser(app.updateLabelHandler))

	router.HandlerFunc(http.MethodPost, "/v1/sublabels/", app.requireActivatedUser(app.createSubLabelHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/sublabels/", app.requireActivatedUser(app.deleteSubLabelHandler))
	// router.HandlerFunc(http.MethodPost, "/v1/blacklist/:id")

	// ? Should path be /v1/labels/... and /v1/sublabels/.... OR /v1/labels/label/... and /v1/labels/sublabel/.....

	router.HandlerFunc(http.MethodPost, "/v1/stories", app.requireActivatedUser(app.createStoryHandler))
	router.HandlerFunc(http.MethodGet, "/v1/stories/:id", app.requireActivatedUser(app.getStoryHandler))
	router.HandlerFunc(http.MethodGet, "/v1/stories", app.requireActivatedUser(app.listStoriesHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/stories/:id", app.requireActivatedUser(app.deleteStoryHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/stories/:id", app.requireActivatedUser(app.updateStoryHandler))

	router.HandlerFunc(http.MethodPost, "/v1/timelines", app.requireActivatedUser(app.createTimelineHandler))
	router.HandlerFunc(http.MethodGet, "/v1/timelines/:id", app.requireActivatedUser(app.getTimelineHandler))
	router.HandlerFunc(http.MethodGet, "/v1/timelines", app.requireActivatedUser(app.listTimelineHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/timelines/:id", app.requireActivatedUser(app.deleteTimelineHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/timelines/:id", app.requireActivatedUser(app.updateTimelineHandler))

	router.HandlerFunc(http.MethodPost, "/v1/events", app.requireActivatedUser(app.createEventHandler))
	router.HandlerFunc(http.MethodGet, "/v1/events/:id", app.requireActivatedUser(app.getEventHandler))
	router.HandlerFunc(http.MethodGet, "/v1/events", app.requireActivatedUser(app.listEventHandler))
	// router.HandlerFunc(http.MethodGet, "/v1/eventsstory", app.requireActivatedUser(app.listStoryEventHandler))
	// ! Change this back
	router.HandlerFunc(http.MethodGet, "/v1/eventsstory/:id", app.requireActivatedUser(app.listStoryEventHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/events/:id", app.requireActivatedUser(app.deleteEventHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/events/:id", app.requireActivatedUser(app.updateEventHandler))

	router.HandlerFunc(http.MethodPost, "/v1/characters", app.requireActivatedUser(app.createCharacterHandler))
	router.HandlerFunc(http.MethodGet, "/v1/characters/:id", app.requireActivatedUser(app.getCharacterHandler))
	router.HandlerFunc(http.MethodGet, "/v1/characters", app.requireActivatedUser(app.listCharacterHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/characters/:id", app.requireActivatedUser(app.deleteCharacterHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/characters/:id", app.requireActivatedUser(app.updateCharacterHandler))

	router.HandlerFunc(http.MethodPost, "/v1/charlabels", app.requireActivatedUser(app.createCharLabelHandler))
	router.HandlerFunc(http.MethodGet, "/v1/charlabels/characters", app.requireActivatedUser(app.listCharacterByLabelsHandler))
	router.HandlerFunc(http.MethodGet, "/v1/charlabels/labels", app.requireActivatedUser(app.listLabelsbyCharacterHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/charlabels", app.requireActivatedUser(app.deleteCharLabelHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/password", app.updateUserPasswordHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", app.createActivationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)
	
	return app.recoverPanic(app.enableCORS(app.authenticate(router)))
}
