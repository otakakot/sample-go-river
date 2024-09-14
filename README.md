# sample-go-river

```shell
go install github.com/riverqueue/river/cmd/river@latest
```

```shell
docker compose up -d postgres
```

```shell
river migrate-up --database-url "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
```

```shell
docker compose up -d
```

```shell
docker exec -it postgres psql -U postgres
```

```shell
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```


```shell
sqlc generate
```
