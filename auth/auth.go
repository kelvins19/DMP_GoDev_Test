package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtSigningSecret = []byte("my_secret")

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Auth struct{}

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

	// Compare password to retrieved password hash
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (auth *Auth) ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSigningSecret, nil
	})
	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp := int64(claims["exp"].(float64))
		iat := int64(claims["iat"].(float64))
		if exp <= time.Now().Unix() {
			return false, fmt.Errorf("token_expired")
		}
		if time.Now().Unix()-iat > 5*60 {
			return false, fmt.Errorf("token_expired")
		}
		return true, nil
	}
	return false, fmt.Errorf("invalid token")
}

// Route handler to handle login requests
func (auth *Auth) LoginHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	// Create JWT with claims
	claims := jwt.MapClaims{}
	claims["sub"] = req.Username
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	claims["iat"] = time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(jwtSigningSecret)
	fmt.Print(err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Error signing token")
		return
	}

	// Return JWT to client as response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}
