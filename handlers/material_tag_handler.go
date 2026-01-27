package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backendLMS/repositories"

	"github.com/gorilla/mux"
)

type attachTagRequest struct {
	TagID int64 `json:"tag_id"`
}

func AttachTagToMaterial(w http.ResponseWriter, r *http.Request) {
	materialID, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	var req attachTagRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.TagID == 0 {
		http.Error(w, "tag_id required", http.StatusBadRequest)
		return
	}

	if err := repositories.AttachTagToMaterial(
		r.Context(),
		materialID,
		req.TagID,
	); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetMaterialTags(w http.ResponseWriter, r *http.Request) {
	materialID, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	data, err := repositories.GetTagsByMaterial(r.Context(), materialID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(data)
}

func DetachMaterialTag(w http.ResponseWriter, r *http.Request) {
	materialID, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	tagID, _ := strconv.ParseInt(mux.Vars(r)["tag_id"], 10, 64)

	repositories.DetachTagFromMaterial(r.Context(), materialID, tagID)
	w.WriteHeader(http.StatusNoContent)
}
