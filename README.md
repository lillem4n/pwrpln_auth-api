## Databaes migration

Done using [dbmate](https://github.com/amacneil/dbmate). Db stuff is stored in `./db`.

Example of running the migrations:

`docker run --rm -it -e DATABASE_URL="postgres://postgres:postgres@127.0.0.1:5432/pwrpln?sslmode=disable" --network=host -v "$(pwd)/db:/db" amacneil/dbmate up`

Example of setting up a postgres SQL server:

`docker run -d --name postgres --network=host -e POSTGRES_PASSWORD=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=pwrlpln postgres`

If you are on of those poor people runnin macOS you must use this one to start a postgres server :(

`docker run -d --name pwrpln-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=pwrlpln postgres`

## Admin account

On first startup with a clean database, an account with name "admin" and the field "role" with a value "admin" is created with no password, using the API Key from ADMIN_API_KEY in the .env file.

## Special account field: "role"

The account field "role" is a bit special, in that if it contains "admin" as one of its values, that grants access to all methods on all accounts on this service. It might be a good idea to use the field "role" for authorization throughout your services.
