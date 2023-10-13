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
	// router.HandlerFunc(http.MethodPost, "/v1/labels/sublabel/:id")
	// router.HandlerFunc(http.MethodPost, "/v1/labels/blacklist/:id")

	return router
}
