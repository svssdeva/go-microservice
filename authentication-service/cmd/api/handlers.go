package main

import (
	"errors"
	"net/http"
)

func(app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	// Implement authentication logic here
	var requestPayload struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("Invalid Credentials"), http.StatusBadRequest)
		return
	}
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	if !valid {
		app.errorJSON(w, errors.New("Invalid Credentials"), http.StatusUnauthorized)
		return
	}
	payload := jsonResponse {
		Error: false,
		Message: "Logged in successfully",
		Data: user,
	}
	
	app.writeJSON(w, http.StatusAccepted, payload)
}
