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

## Endpoints
- GET /v1/healthcheck - Check the health of the API
- POST /v1/movies - Create a new movie
- GET /v1/movies/:id - Get a specific movie by ID
- PATCH /v1/movies/:id - Update a specific movie by ID
- DELETE /v1/movies/:id - Delete a specific movie by ID