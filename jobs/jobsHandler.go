package jobs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/kelvins19/DMP_Test/auth"
)

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

type Jobs struct{}

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

// Function to retrieve detail of job from external endpoint
func getJobDetailFromEndpoint(id string) (jobData, error) {
	// Make HTTP request to external endpoint to retrieve job detail
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
func (jobs *Jobs) JobListHandler(w http.ResponseWriter, r *http.Request) {
	// Verify JWT in request header
	tokenString := r.Header.Get("Authorization")

	auth := &auth.Auth{}
	valid, err := auth.ValidateToken(tokenString)
	if !valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
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
func (jobs *Jobs) JobDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Verify JWT in request header
	tokenString := r.Header.Get("Authorization")

	auth := &auth.Auth{}
	valid, err := auth.ValidateToken(tokenString)
	if !valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	id := r.URL.Query().Get("id")

	// Retrieve job detail from external endpoint
	jobDetail, err := getJobDetailFromEndpoint(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return job detail to client as response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobDetail)
}
