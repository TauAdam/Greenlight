package main

import (
	"fmt"
	"golang.org/x/time/rate"
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

// global rate limiter to enforce strict limit on the total rate of requests to API
func (app *application) rateLimit(next http.Handler) http.Handler {
	//  allows an average of 2 requests per second, with a maximum of 4 requests in a single ‘burst’.
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if the request not permitted by the limiter, send a 429 Too Many Requests response.
		if !limiter.Allow() {
			app.rateLimitExceededResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
