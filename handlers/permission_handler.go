package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backendLMS/models"
	"backendLMS/repositories"

	"github.com/gorilla/mux"
)

func CreatePermission(w http.ResponseWriter, r *http.Request) {
	var req models.Permission
	json.NewDecoder(r.Body).Decode(&req)

	if req.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	if err := repositories.CreatePermission(r.Context(), &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}

func GetPermissions(w http.ResponseWriter, r *http.Request) {
	data, _ := repositories.GetAllPermissions(r.Context())
	json.NewEncoder(w).Encode(data)
}

func GetPermissionByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	data, err := repositories.GetPermissionByID(r.Context(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func UpdatePermission(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	var req models.Permission
	json.NewDecoder(r.Body).Decode(&req)

	if err := repositories.UpdatePermission(r.Context(), id, &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeletePermission(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	repositories.DeletePermission(r.Context(), id)
	w.WriteHeader(http.StatusNoContent)
}
