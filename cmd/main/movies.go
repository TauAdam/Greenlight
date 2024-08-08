package main

import (
	"fmt"
	"net/http"
)

func (app *application) handleCreateMovie(w http.ResponseWriter, r *http.Request) {
	// TODO - create a new movie
	_, err := fmt.Fprintln(w, "create a new movie")
	if err != nil {
		return
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