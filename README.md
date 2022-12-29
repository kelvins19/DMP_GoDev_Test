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
4. Make sure you have migrate the database before go to the next step
5. Run `go run main.go` to start the application

## Test Login Credential
```
Username: kelvins19
Password: 123456
```

## Endpoints List
- POST localhost:8080/login
- GET localhost:8080/jobs?page=&description=&location=&full_time=
- GET localhost:8080/job-detai?id={job-id}

## Task Details
1. Login API (POST localhost:8080/login)
    - The API should validate username and password
    - List of valid username and password should be stored on a DBMS 
    - Any DBMS is allowed
    - The API should implement JSON Web Token (JWT)


2. Get job list API (GET localhost:8080/jobs?page=&description=&location=&full_time=)
    - The API should be secured with JWT Authorization
    -  The API should make http request to
http://dev3.dansmultipro.co.id/api/recruitment/positions.json and return
jobs data as response payload.
    - The API should provide “search” functionality tearch for jobs by term,
location, full time vs part time, or any combination of the three. All parameters are optional.
        - **description** — A search term, such as "ruby" or "java". This parameter is aliased to search.
        - **location** — A city name, zip code, or other location search term.
        - **full_time** — If you want to limit results to full time positions set this
parameter to 'true'.
    - The API should support pagination Example
http://dev3.dansmultipro.co.id/api/recruitment/positions.json?page=1


3. Get job detail API (GET localhost:8080/job-detai?id={job-id})
    - The API should be secured with JWT Authorization
    - The API should make http request
http://dev3.dansmultipro.co.id/api/recruitment/positions/{ID} and return job detail data as response payload.