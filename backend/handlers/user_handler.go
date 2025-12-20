package handlers

import (
	"encoding/json"
	"net/http"

	"backendLMS/middlewares"
)

func Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	roleID := r.Context().Value(middlewares.CtxRoleID).(int64)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userID,
		"role_id": roleID,
	})
}