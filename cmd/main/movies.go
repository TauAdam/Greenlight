package main

import (
	"fmt"
	"net/http"
)

func (app *application) handleCreateMovie(w http.ResponseWriter, r *http.Request) {
	var input struct {

Title string `json:"title"`

Year int32 `json:"year"`
Runtime int32 `json:"runtime"`
Genres []string `json:"genres"`
}
// Use the new readJSON() helper to decode the request body into the input struct.
err := app.readJSON(w, r, &input)
if err != nil {
app.errorResponse(w, r, http.StatusBadRequest, err.Error())
return
}
fmt.Fprintf(w, "%+v\n", input)
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