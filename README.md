A json api for managing and retrieving information about movies. 

## Self-host instructions 
Prerequisites
- Make sure you have docker, docker compose plugin installed on your machine


- Download the project
- Build and run the docker container
   ```shell
   docker compose up --build
   ```
- run ```db/migrations/up``` to apply migrations to the database
- the api ready to use on `http://localhost:4000`

![img.png](img.png)
![img_2.png](img_2.png)

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
- Middleware for ip based rate limiting
    - separate rate limiter for each client, so that one bad
      client making too many requests doesn’t affect all the others.
- Configuration for
    - Database
    - Server
    - Rate limiter
- Graceful shutdown
    - The server will wait for all active connections to finish before shutting down.
- Sending mail
- User Activation
    - When a user registers, they will receive an email with a link to activate their account.
- Token based authentication
- CORS

## Endpoints
- GET /v1/healthcheck - Check the health of the API
- POST /v1/movies - Create a new movie
- GET /v1/movies/:id - Get a specific movie by ID
- PATCH /v1/movies/:id - Update a specific movie by ID
- DELETE /v1/movies/:id - Delete a specific movie by ID

![img_1.png](img_1.png)