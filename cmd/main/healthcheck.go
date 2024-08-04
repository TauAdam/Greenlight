package main

import (
	"fmt"
	"net/http"
)

func (app *application) handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "status: available")
	if err != nil {
		app.logger.Print(w, err)
		return
	}
	_, err = fmt.Fprintf(w, "environment: %s\n", app.config.env)
	if err != nil {
		app.logger.Print(w, err)
		return
	}
	_, err = fmt.Fprintf(w, "version: %s\n", version)
	if err != nil {
		app.logger.Print(w, err)
		return
	}
}
