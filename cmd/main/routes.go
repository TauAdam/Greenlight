package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// custom 404, not allowed method handler
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/v1/users", app.handleRegisterUser)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.handleActivateUser)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.handleCreateAuthenticationToken)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.handleHealthcheck)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.handleCreateMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.handleShowMovie)
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.handleUpdateMovie)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.handleDeleteMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.handleListMovies)

	return app.recoverPanic(app.rateLimit(router))
}
