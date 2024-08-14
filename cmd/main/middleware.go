package main

import (
	"fmt"
	"net/http"
)

// middleware to recover from panics & log them
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// deferred function (which will always be run in the event of a panic as Go unwinds the stack).
		defer func() {
			// check if there has been a panic or not.
			if err := recover(); err != nil {
				// This acts as a trigger to make Go's HTTP server
				// automatically close the current connection after a response has been sent.
				w.Header().Set("Connection", "close")
				//  use fmt.Errorf() to normalize any type into an error & log at the ERROR level
				// & send a 500 Internal Server Error response.
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
