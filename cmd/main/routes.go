package main

import (
	"expvar"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	PermissionRead  = "movies:read"
	PermissionWrite = "movies:write"
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

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission(PermissionRead, app.handleListMovies))
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission(PermissionWrite, app.handleCreateMovie))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission(PermissionRead, app.handleShowMovie))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission(PermissionWrite, app.handleUpdateMovie))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission(PermissionWrite, app.handleDeleteMovie))

	router.Handler(http.MethodGet, "/v1/metrics", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
