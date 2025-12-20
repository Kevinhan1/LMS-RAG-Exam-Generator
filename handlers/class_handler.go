package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"backendLMS/models"
	"backendLMS/repositories"

	"github.com/gorilla/mux"
)

func CreateClass(w http.ResponseWriter, r *http.Request) {
	var c models.Class

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// VALIDATION
	if err := validateClass(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now().Unix()
	c.TimeCreated = now
	c.TimeModified = now

	if err := repositories.CreateClass(r.Context(), &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func GetClasses(w http.ResponseWriter, r *http.Request) {
	data, _ := repositories.GetClasses(context.Background())
	json.NewEncoder(w).Encode(data)
}

// GET /api/classes/{id}
func GetClassByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	class, err := repositories.GetClassByID(r.Context(), id)
	if err != nil {
		http.Error(w, "class not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(class)
}

// PUT /api/admin/classes/{id}
func UpdateClass(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var c models.Class
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// VALIDATION
	if err := validateClass(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.ID = id
	c.TimeModified = time.Now().Unix()

	if err := repositories.UpdateClass(r.Context(), &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(c)
}

// DELETE /api/admin/classes/{id}
func DeleteClass(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	if err := repositories.DeleteClass(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func validateClass(c *models.Class) error {
	if c.Grade <= 0 {
		return errors.New("grade must be greater than 0")
	}

	validJenjang := map[string]bool{
		"SD": true, "SMP": true, "SMA": true,
	}

	if !validJenjang[c.Jenjang] {
		return errors.New("invalid jenjang")
	}
	return nil
}
