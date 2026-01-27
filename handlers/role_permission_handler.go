package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backendLMS/repositories"

	"github.com/gorilla/mux"
)

type assignPermissionRequest struct {
	PermissionID int64 `json:"permission_id"`
}

func AssignPermission(w http.ResponseWriter, r *http.Request) {
	roleID, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	var req assignPermissionRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.PermissionID == 0 {
		http.Error(w, "permission_id required", http.StatusBadRequest)
		return
	}

	if err := repositories.AssignPermissionToRole(
		r.Context(),
		roleID,
		req.PermissionID,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetRolePermissions(w http.ResponseWriter, r *http.Request) {
	roleID, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	data, err := repositories.GetPermissionsByRole(r.Context(), roleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func RemovePermission(w http.ResponseWriter, r *http.Request) {
	roleID, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	permID, _ := strconv.ParseInt(mux.Vars(r)["permission_id"], 10, 64)

	repositories.RemovePermissionFromRole(r.Context(), roleID, permID)
	w.WriteHeader(http.StatusNoContent)
}
