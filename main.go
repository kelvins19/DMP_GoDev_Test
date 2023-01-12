package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/kelvins19/DMP_Test/auth"
	"github.com/kelvins19/DMP_Test/jobs"
	_ "github.com/lib/pq"
)

func main() {
	// Connect to PostgreSQL database
	db, err := sql.Open("postgres", "postgres://kelvins19:123456@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	auth := &auth.Auth{}
	jobs := &jobs.Jobs{}

	// Set up route to handle login requests
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		auth.LoginHandler(db, w, r)
	})
	http.HandleFunc("/jobs", jobs.JobListHandler)
	http.HandleFunc("/job-detail", func(w http.ResponseWriter, r *http.Request) {
		jobs.JobDetailHandler(w, r)
	})

	// Start server and listen for requests
	http.ListenAndServe(":8080", nil)
}
