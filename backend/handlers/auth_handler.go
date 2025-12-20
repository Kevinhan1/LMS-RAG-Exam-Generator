package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"backendLMS/models"
	"backendLMS/repositories"
	"backendLMS/services"

	"golang.org/x/crypto/bcrypt"
)

/* ================= REGISTER ================= */

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int64  `json:"role_id"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		http.Error(w, "failed hash password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		RoleID:       req.RoleID,
		TimeCreated:  time.Now().Unix(),
	}

	if err := repositories.CreateUser(context.Background(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

/* ================= LOGIN ================= */

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := repositories.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		http.Error(w, "email atau password salah", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		http.Error(w, "email atau password salah", http.StatusUnauthorized)
		return
	}

	token, err := services.GenerateJWT(user.ID, user.RoleID)
	if err != nil {
		http.Error(w, "failed generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(loginResponse{
		Token: token,
	})
}