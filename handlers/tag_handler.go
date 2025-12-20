package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"backendLMS/middlewares"
	"backendLMS/models"
	"backendLMS/repositories"

	"github.com/gorilla/mux"
)

func CreateTagHandler(w http.ResponseWriter, r *http.Request) {
	var t models.Tag
	json.NewDecoder(r.Body).Decode(&t)

	err := repositories.CreateTag(context.Background(), &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// log activity
	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	repositories.CreateLog(context.Background(), &models.LogActivity{
		UserID:      userID,
		Action:      "create_tag",
		TargetTable: "tags",
		TargetID:    t.ID,
		Description: t.Name,
		CreatedAt:   time.Now(),
	})

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)
}

func GetTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, _ := repositories.GetTags(context.Background())
	json.NewEncoder(w).Encode(tags)
}

func UpdateTagHandler(w http.ResponseWriter, r *http.Request) {
	var t models.Tag
	json.NewDecoder(r.Body).Decode(&t)

	err := repositories.UpdateTag(context.Background(), &t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// log activity
	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	repositories.CreateLog(context.Background(), &models.LogActivity{
		UserID:      userID,
		Action:      "update_tag",
		TargetTable: "tags",
		TargetID:    t.ID,
		Description: t.Name,
		CreatedAt:   time.Now(),
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(t)
}

func DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	err := repositories.DeleteTag(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// log activity
	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	repositories.CreateLog(context.Background(), &models.LogActivity{
		UserID:      userID,
		Action:      "delete_tag",
		TargetTable: "tags",
		TargetID:    id,
		CreatedAt:   time.Now(),
	})

	w.WriteHeader(http.StatusNoContent)
}

func GetTagByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	tag, err := repositories.GetTagByID(context.Background(), id)
	if err != nil {
		http.Error(w, "tag not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(tag)
}
