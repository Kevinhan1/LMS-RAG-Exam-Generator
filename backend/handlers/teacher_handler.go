package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"backendLMS/models"
	"backendLMS/repositories"

	"golang.org/x/crypto/bcrypt"
)

type registerTeacherRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	NIP      string `json:"nip"`
}

func RegisterTeacher(w http.ResponseWriter, r *http.Request) {
	var req registerTeacherRequest
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
		RoleID:       2, // TEACHER
		NIP:          &req.NIP,
		TimeCreated:  now,
	}

	if err := repositories.CreateUser(r.Context(), &user); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
