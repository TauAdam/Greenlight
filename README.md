# Greenlight
A json api for managing and retrieving information about movies. 

## Stack
- Pure Go
- Postgresql

## Features
- Create a movie
- Get a movie
- Get all movies
- Update a movie
- Delete a movie
- Healthcheck
- Filtering, sorting and pagination
- Custom json logging
- Middleware for panic recovery
    - if there is a panic in one of API handlers the middleware
      will recover it and log the error using custom Logger and send the client a nice
      500 Internal Server Error response with a JSON body.
- Middleware for rate limiting
    - The rate limiter middleware is used to limit the number of requests a client can make to the API within a given time window.

## Endpoints
- GET /v1/healthcheck - Check the health of the API
- POST /v1/movies - Create a new movie
- GET /v1/movies/:id - Get a specific movie by ID
- PATCH /v1/movies/:id - Update a specific movie by ID
- DELETE /v1/movies/:id - Delete a specific movie by ID