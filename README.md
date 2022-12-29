## README

## Tech Stack
- Golang V1.19
- PostgreSQL as DBMS

## How to migrate
1. Install golang-migrate[https://dev.to/techschoolguru/how-to-write-run-database-migration-in-golang-5h6g]
2. Run migration
> migrate -path db/migration -database "postgresql://username:password@localhost:5432/postgres?sslmode=disable" -verbose up

## How to setup the application
1. Go to the project directory
2. Run `go build`
3. Run `go mod tidy`
4. Run `go run main.go` to start the application

## Test Login Credential
```
Username: kelvins19
Password: 123456
```

## Endpoints List
- POST localhost:8080/login
- GET localhost:8080/jobs?page=&description=&location=&full_time=
- GET localhost:8080/job-detai?id={job-id}