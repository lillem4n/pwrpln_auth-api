# Auth API

A tiny REST API for auth. Register accounts, auth with api-key or name/password, renew JWT tokens...

## Quick start with docker compose

Migrate database: `docker-compose run --rm db-migrations`

Start the API (on port 4000 by default): `docker-compose up -d`

Point your browser to `http://localhost:4000` to view the swagger API documentation.

## Admin account

On first startup with a clean database, an account with name "admin" and the field "role" with a value "admin" is created with no password, using the API Key from ADMIN_API_KEY in the .env file.

## Special account field: "role"

The account field "role" is a bit special, in that if it contains "admin" as one of its values, that grants access to all methods on all accounts on this service. It might be a good idea to use the field "role" for authorization throughout your services.

## Tests

Run integration tests (Requires migrated database and started API): `docker-compose run --rm tests`

## Some useful cURLs

Obtain an admin GWT: `curl -d '"api-key-goes-here"' -H "Content-Type: application/json" -i http://localhost:4000/auth/api-key`

Use a bearer token to make a call: `curl -H "Content-Type: application/json" -H "Authorization: bearer your-JWT-token-goes-here" -i http://localhost:4000/account/{accountID}`

Create account: `curl -d '{"name": "Bosse", "password": "Hemligt", "fields": [{ "name":"role", "values":["user"]}]}' -H "Content-Type: application/json" -H "Authorization: bearer your-JWT-token-goes-here" -i http://localhost:4000/account`

x
