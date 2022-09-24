package app

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

type requestCredentials struct {
	UserName string `json:"userName" validate:"required,alphanum,gte=4"`
	Password string `json:"password" validate:"required,gte=8"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var credentials requestCredentials

	// decode the request body
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate credentials
	credError := validateCredentials(credentials)
	if credError != nil {
		http.Error(w, credError.Error(), http.StatusBadRequest)
		return
	}

	// generate id
	id := uuid.New().String()

	// add new user
	users[credentials.UserName] = &User{id, credentials.UserName, credentials.Password}

	// create 200 response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseData := map[string]string{"id": id}
	json.NewEncoder(w).Encode(responseData)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials requestCredentials

	// decode the request body, error returns 400
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if user exists and get id, error returns 400
	id, ok := getUserId(credentials, users)
	if !ok {
		http.Error(w, "Invalid username/Password", http.StatusBadRequest)
		return
	}

	// generate token, error returns 500
	token, expiring, err := generateToken(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create 200 response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("X-Rate-Limit", "1")
	w.Header().Set("X-Expires-After", expiring)
	responseData := map[string]string{"url": "ws://fancy-chat.io/ws&token=" + token}
	json.NewEncoder(w).Encode(responseData)
}
