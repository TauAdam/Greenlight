package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
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

// ip based rate limiter for each client
func (app *application) rateLimit(next http.Handler) http.Handler {
	var (
		mu      sync.Mutex
		clients = make(map[string]*rate.Limiter)
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Lock the mutex to prevent concurrent execution
		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = rate.NewLimiter(2, 4)
		}

		if !clients[ip].Allow() {
			// Unlock the mutex before sending a response
			mu.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}

		// Unlock the mutex before calling the next handler in the chain
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
