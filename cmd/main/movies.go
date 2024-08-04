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
	id, err := app.extractIDParam(r)
	if err != nil {
		app.logger.Printf("invalid movie ID %q", id)
		http.NotFound(w, r)
		return
	}

	_, err = fmt.Fprintf(w, "show the details of movie %d\n", id)
	if err != nil {
		app.logger.Print(w, err)
		return
	}
}
