package main

import (
	"errors"
	"github.com/TauAdam/Greenlight/internal/data"
	"github.com/TauAdam/Greenlight/internal/validator"
	"net/http"
)

func (app *application) handleCreateAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePassword(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusUnauthorized, "invalid credentials")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

}
