package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/labels/:id", app.showLabelHandler)
	router.HandlerFunc(http.MethodPost, "/v1/labels", app.createLabelHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/labels/:id", app.deleteLabelHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/labels/:id", app.updateLabelHandler)
	// router.HandlerFunc(http.MethodPost, "/v1/sublabel/:id")
	// router.HandlerFunc(http.MethodPost, "/v1/blacklist/:id")

	// ? Should path be /v1/labels/... and /v1/sublabels/.... OR /v1/labels/label/... and /v1/labels/sublabel/.....

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	return router
}
