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
	// router.HandlerFunc(http.MethodPost, "/v1/sublabel/:id")
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
	router.HandlerFunc(http.MethodDelete, "/v1/events/:id", app.requireActivatedUser(app.deleteEventHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/events/:id", app.requireActivatedUser(app.updateEventHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthTokenHandler)
	return app.authenticate(router)
}
