package main

import (
	"errors"
	"fmt"
	"github.com/TauAdam/Greenlight/internal/data"
	"github.com/TauAdam/Greenlight/internal/validator"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"strings"
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

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()
		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}

// guard that check that a user is both authenticated and activated.
func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	activated := func(w http.ResponseWriter, r *http.Request) {
		user := app.contextRetrieveUser(r)

		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}

	return app.requireAuthenticatedUser(activated)
}

func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := app.contextRetrieveUser(r)
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// guard that require permission, activation
func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	permitted := func(w http.ResponseWriter, r *http.Request) {
		user := app.contextRetrieveUser(r)

		permissions, err := app.models.Permissions.PermissionsOfUser(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		if !permissions.Includes(code) {
			app.notPermittedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}

	return app.requireActivatedUser(permitted)
}

// enableCORS middleware to add the necessary CORS headers to the response
func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")

		origin := r.Header.Get("Origin")

		// If the Origin header is empty, the request is same-origin, and we don't need to do anything.
		if origin != "" && len(app.config.cors.trustedOrigins) != 0 {
			// Loop through the list of trusted origins
			for i := range app.config.cors.trustedOrigins {
				// If we find a match, set the Access-Control-Allow-Origin header on the response
				if origin == app.config.cors.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
