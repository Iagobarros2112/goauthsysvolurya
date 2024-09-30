package main

//In healthcheckHandler, we also returned a
//server error in case there was an error "serializing"
//the data.
//This error is defined in cmd/api/errors.go:
//serverErrorResponse uses errorResponse
//to tell the user about the unexpected
//issue with our app. Both of them used
//logError to log the error to our logging
//console so that we can easily debug it.
//By now, you should see how beautiful
//it is with the "polymorphism" we have
//implemented! We can just use a method
//anywhere without breaking a sweat!

import (
	"net/http"
)

type envelope map[string]interface{}

func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}
