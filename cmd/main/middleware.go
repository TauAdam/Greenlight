package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
	"time"
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
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)
	// background goroutine which removes old entries from the clients map once every minute.
	go func() {
		for {
			time.Sleep(time.Minute)
			// Lock the mutex to prevent any other goroutines from accessing the clients map while the cleanup is taking place.
			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			// unlock the mutex when the cleanup is complete.
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}
			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})
}
