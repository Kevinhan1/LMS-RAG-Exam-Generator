package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"backendLMS/db"
	"backendLMS/models"

	"github.com/gorilla/mux"
)

// CreateRole - POST /roles
func CreateRole(w http.ResponseWriter, r *http.Request) {
	var input models.Role
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	now := time.Now().Unix()
	sql := `INSERT INTO roles (name, description, timecreated, timemodified) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int64
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.Pool.QueryRow(ctx, sql, input.Name, input.Description, now, now).Scan(&id)
	if err != nil {
		http.Error(w, "failed to insert role: "+err.Error(), http.StatusInternalServerError)
		return
	}
	input.ID = id
	input.TimeCreated = now
	input.TimeModified = now

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

// GetRoles - GET /roles
func GetRoles(w http.ResponseWriter, r *http.Request) {
	sql := `SELECT id, name, description, timecreated, timemodified FROM roles ORDER BY id`
	rows, err := db.Pool.Query(context.Background(), sql)
	if err != nil {
		http.Error(w, "failed to query roles: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	roles := []models.Role{}
	for rows.Next() {
		var rmodel models.Role
		if err := rows.Scan(&rmodel.ID, &rmodel.Name, &rmodel.Description, &rmodel.TimeCreated, &rmodel.TimeModified); err != nil {
			http.Error(w, "failed to scan row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		roles = append(roles, rmodel)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

// GetRole - GET /roles/{id}
func GetRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	sql := `SELECT id, name, description, timecreated, timemodified FROM roles WHERE id=$1`
	var rmodel models.Role
	err = db.Pool.QueryRow(context.Background(), sql, id).Scan(&rmodel.ID, &rmodel.Name, &rmodel.Description, &rmodel.TimeCreated, &rmodel.TimeModified)
	if err != nil {
		http.Error(w, "role not found: "+err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rmodel)
}

// UpdateRole - PUT /roles/{id}
func UpdateRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var input models.Role
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	now := time.Now().Unix()
	sql := `UPDATE roles SET name=$1, description=$2, timemodified=$3 WHERE id=$4`
	ct, err := db.Pool.Exec(context.Background(), sql, input.Name, input.Description, now, id)
	if err != nil {
		http.Error(w, "failed to update role: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if ct.RowsAffected() == 0 {
		http.Error(w, "role not found", http.StatusNotFound)
		return
	}

	input.ID = id
	input.TimeModified = now
	// note: we don't set TimeCreated here; client can fetch GET /roles/{id} if needed

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(input)
}

// DeleteRole - DELETE /roles/{id}
func DeleteRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	sql := `DELETE FROM roles WHERE id=$1`
	ct, err := db.Pool.Exec(context.Background(), sql, id)
	if err != nil {
		http.Error(w, "failed to delete role: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if ct.RowsAffected() == 0 {
		http.Error(w, "role not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
