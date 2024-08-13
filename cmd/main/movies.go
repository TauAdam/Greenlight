package main

import (
	"errors"
	"fmt"
	"github.com/TauAdam/Greenlight/internal/data"
	"github.com/TauAdam/Greenlight/internal/validator"
	"log"
	"net/http"
	"strconv"
)

func (app *application) handleCreateMovie(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`

		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	// Use the new readJSON() helper to decode the request body into the input struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// validation step

	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// include a Location header to let the
	// client know which URL they can find the newly-created resource at
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	// Write a JSON response with a 201 code
	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// handleShowMovie retrieve the interpolated "id" parameter from the current URL
func (app *application) handleShowMovie(w http.ResponseWriter, r *http.Request) {
	id, err := app.extractIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) handleUpdateMovie(w http.ResponseWriter, r *http.Request) {
	id, err := app.extractIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.FormatInt(int64(movie.Version), 32) != r.Header.Get("X-Expected-Version") {
			app.editConflictResponse(w, r)
			return
		}
	}

	var input struct {
		Title   *string       `json:"title"`
		Year    *int32        `json:"year"`
		Runtime *data.Runtime `json:"runtime"`
		Genres  []string      `json:"genres"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) handleDeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := app.extractIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// handleListMovies implements filtering, sorting, pagination
func (app *application) handleListMovies(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}

	v := validator.New()

	paramsMap := r.URL.Query()
	input.Title = app.readString(paramsMap, "title", "")
	input.Genres = app.readCSV(paramsMap, "genres", []string{})

	input.Filters.Page = app.readInteger(paramsMap, "page", 1, v)
	input.Filters.PageSize = app.readInteger(paramsMap, "page_size", 20, v)
	input.Filters.Sort = app.readString(paramsMap, "sort", "id")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err := app.writeJSON(w, http.StatusAccepted, envelope{"parsed-query-params:": input}, nil)
	if err != nil {
		log.Printf("failed to send response")
		return
	}
}
