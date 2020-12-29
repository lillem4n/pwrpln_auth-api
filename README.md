## Databaes migration

Done using [dbmate](https://github.com/amacneil/dbmate). Db stuff is stored in `./db`.

Example of running the migrations:

`docker run --rm -it -e DATABASE_URL="postgres://postgres:postgres@127.0.0.1:5432/pwrpln?sslmode=disable" --network=host -v "$(pwd)/db:/db" amacneil/dbmate up`

Example of setting up a postgres SQL server:

`docker run -d --name postgres --network=host -e POSTGRES_PASSWORD=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=pwrlpln postgres`