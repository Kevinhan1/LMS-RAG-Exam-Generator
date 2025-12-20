package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"backendLMS/repositories"
	"github.com/gorilla/mux"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repositories.GetAllUsers(context.Background())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	var body struct {
		RoleID int64 `json:"role_id"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	err := repositories.UpdateUserRole(context.Background(), id, body.RoleID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	err := repositories.DeleteUser(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}