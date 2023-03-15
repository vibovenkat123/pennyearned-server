package apiHelpers

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"errors"
)

var (
	ErrInvalidID = errors.New("The ID provided is not the valid format for a id")
	ErrExpenseInvalid = errors.New("The parameters provided for creating/updating a expense are invalid")
)

// basic error log
func (app *Application) LogError(err error, r *http.Request) {
	App.Log.Error("Error in the API",
		zap.Error(err),
	)
}

// the simplest error helper
func (app *Application) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := Envelope{"error": message}
	err := App.WriteJSON(w, status, env, nil)
	if err != nil {
		App.LogError(err, r)
		w.WriteHeader(500)
	}
}

// the server error (500)
func (app *Application) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	App.LogError(err, r)
	message := "the server encountered a problem and could not process your request"
	App.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

// a 404 not found error

func (app *Application) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	App.ErrorResponse(w, r, http.StatusNotFound, message)
}

func (app *Application) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	App.ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *Application) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	App.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *Application) ConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "The request conflicted with the current state of the server"
	App.ErrorResponse(w, r, http.StatusConflict, message)
}
