package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"backendLMS/models"
	"backendLMS/repositories"

	"golang.org/x/crypto/bcrypt"
)

type registerStudentRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	NISN     string `json:"nisn"`
	ClassID  int64  `json:"class_id"`
}

func RegisterStudent(w http.ResponseWriter, r *http.Request) {
	var req registerStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	now := time.Now().Unix()

	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		RoleID:       3, // STUDENT
		ClassID:      &req.ClassID,
		NISN:         &req.NISN,
		TimeCreated:  now,
	}

	if err := repositories.CreateUser(r.Context(), &user); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
