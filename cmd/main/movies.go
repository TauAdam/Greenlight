package main

import (
	"fmt"
	"net/http"
)

func (app *application) handleCreateMovie(w http.ResponseWriter, r *http.Request) {
	var input struct {

Title string `json:"title"`

Year int32 `json:"year"`
Runtime data.Runtime `json:"runtime"`
Genres []string `json:"genres"`
}
// Use the new readJSON() helper to decode the request body into the input struct.
err := app.readJSON(w, r, &input)
if err != nil {
app.badRequestResponse(w, r, err)
return
}
// validation step

movie := &data.Movie{
Title: input.Title,
Year: input.Year,
Runtime: input.Runtime,
Genres: input.Genres,
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
id, err := app.readIDParam(r)
if err != nil {
http.NotFound(w, r)
return
}
movie := data.Movie{
ID: id,
CreatedAt: time.Now(),
Title: "Casablanca",
Runtime: 102,
Genres: []string{"drama", "romance", "war"},
Version: 1,
}

err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
if err != nil {
   app.serverErrorResponse(w, r, err)
}
}