package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const jwtSigningSecret = "my_secret"

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type jobData struct {
	Id          string `json:"id"`
	Type        string `json:"type"`
	Url         string `json:"url"`
	CreatedAt   string `json:"created_at"`
	Company     string `json:"company"`
	CompanyUrl  string `json:"company_url"`
	Location    string `json:"location"`
	Title       string `json:"title"`
	Description string `json:"description"`
	HowToApply  string `json:"how_to_apply"`
	CompanyLogo string `json:"company_logo"`
}

// Function to authenticate user based on provided username and password
func authenticate(db *sql.DB, username, password string) bool {
	// Retrieve password hash for provided username from database
	var hash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE username = $1", username).Scan(&hash)
	if err == sql.ErrNoRows {
		// Return false if no matching user is found
		return false
	} else if err != nil {
		// Return false if there is an error executing the query
		return false
	}

	// Compare provided password to retrieved password hash
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Route handler to handle login requests
func loginHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Parse login request from request body
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Authenticate user based on provided username and password
	if !authenticate(db, req.Username, req.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Generate JWT for authenticated user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": req.Username,
	})
	tokenString, err := token.SignedString([]byte(jwtSigningSecret))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return JWT to client as response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

// Function to retrieve list of jobs from external endpoint
func getJobListFromEndpoint(page, description, location, fullTime string) ([]jobData, error) {
	// Make HTTP request to external endpoint to retrieve job list
	url := fmt.Sprintf("http://dev3.dansmultipro.co.id/api/recruitment/positions.json?page=%s&description=%s&location=%s&full_time=%s", page, url.QueryEscape(description), url.QueryEscape(location), url.QueryEscape(fullTime))
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("1")
		return nil, err
	}
	defer resp.Body.Close()

	// Parse job list from response
	var jobList []jobData
	err = json.NewDecoder(resp.Body).Decode(&jobList)
	if err != nil {
		return nil, err
	}

	return jobList, nil
}

// Function to retrieve detail of jobs from external endpoint
func getJobDetailFromEndpoint(id string) (jobData, error) {
	// Make HTTP request to external endpoint to retrieve job list
	var jobDetail jobData
	resp, err := http.Get("http://dev3.dansmultipro.co.id/api/recruitment/positions/" + id)
	if err != nil {
		return jobDetail, err
	}
	defer resp.Body.Close()

	// Parse job detail from response
	err = json.NewDecoder(resp.Body).Decode(&jobDetail)
	if err != nil {
		return jobDetail, err
	}

	return jobDetail, nil
}

// Route handler to handle requests to retrieve job list
func jobListHandler(w http.ResponseWriter, r *http.Request) {
	// Verify JWT in request header
	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSigningSecret), nil
	})
	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	page := r.URL.Query().Get("page")
	description := r.URL.Query().Get("description")
	location := r.URL.Query().Get("location")
	fullTime := r.URL.Query().Get("full_time")

	if fullTime == "" {
		fullTime = "false"
	}

	// Retrieve job list from external endpoint
	jobList, err := getJobListFromEndpoint(page, description, location, fullTime)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return job list to client as response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobList)
}

// Route handler to handle requests to retrieve job detail based on id
func jobDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Verify JWT in request header
	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSigningSecret), nil
	})
	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	id := r.URL.Query().Get("id")

	// Retrieve job list from external endpoint
	jobDetail, err := getJobDetailFromEndpoint(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return job detail to client as response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobDetail)
}

func main() {
	// Connect to PostgreSQL database
	db, err := sql.Open("postgres", "postgres://kelvins19:123456@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Set up route to handle login requests
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginHandler(db, w, r)
	})
	http.HandleFunc("/jobs", jobListHandler)
	http.HandleFunc("/job-detail", func(w http.ResponseWriter, r *http.Request) {
		jobDetailHandler(w, r)
	})

	// Start server and listen for requests
	http.ListenAndServe(":8080", nil)
}
