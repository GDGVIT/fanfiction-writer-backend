package main

import (
	"fmt"
	"net/http"
)

// logError is a generic helper to log an error message
func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// errorResponse is a generic helper for sending JSON-formatted error messages to the client
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// serverErrorResponse is used when the server encounters an error during runtime. It logs the detailed error, then uses the errorResponse() helper to send a 500 Internal Server Error status code and a generic error message
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "The server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// notFoundResponse is used to send a 404 Not Found status code and an appropriate error message
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowedResponse is used to send a 405 Method Not Allowed status code and an appropriate error message
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not allowed for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// badRequestResponse is used to send a 400 Bad Request status code and an appropriate error message
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// failedValidationResponse is used to send a 422 Unprocessable Entity status code and an appropriate error message
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// editConflictResponse is used to send a 409 Status conflict status code and an appropriate error message
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

// invalidCredentialsResponse is used to send a 401 status unauthorized status code and an appropriate error message
func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request)  {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// invalidAuthenticationTokenResponse is used to send a 401 status unauthorized status code and an appropriate error message
func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request)  {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request)  {
	message := "your account must be activated to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

