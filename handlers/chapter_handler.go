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

func CreateChapterHandler(w http.ResponseWriter, r *http.Request) {
	var c models.Chapter
	json.NewDecoder(r.Body).Decode(&c)

	err := repositories.CreateChapter(context.Background(), &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// log activity
	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	repositories.CreateLog(context.Background(), &models.LogActivity{
		UserID:      userID,
		Action:      "create_chapter",
		TargetTable: "chapters",
		TargetID:    c.ID,
		Description: c.Title,
		CreatedAt:   time.Now(),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func GetChaptersHandler(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.ParseInt(mux.Vars(r)["course_id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid course id", http.StatusBadRequest)
		return
	}

	chapters, err := repositories.GetChapters(context.Background(), courseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(chapters)
}

func UpdateChapterHandler(w http.ResponseWriter, r *http.Request) {
	var c models.Chapter

	// ambil ID dari URL (source of truth)
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid chapter id", http.StatusBadRequest)
		return
	}

	// decode body
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// override ID dari URL (ANTI TAMPERING)
	c.ID = id

	// update
	if err := repositories.UpdateChapter(context.Background(), &c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// log activity
	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	repositories.CreateLog(context.Background(), &models.LogActivity{
		UserID:      userID,
		Action:      "update_chapter",
		TargetTable: "chapters",
		TargetID:    c.ID,
		Description: c.Title,
		CreatedAt:   time.Now(),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}

func DeleteChapterHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	err := repositories.DeleteChapter(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// log activity
	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	repositories.CreateLog(context.Background(), &models.LogActivity{
		UserID:      userID,
		Action:      "delete_chapter",
		TargetTable: "chapters",
		TargetID:    id,
		CreatedAt:   time.Now(),
	})

	w.WriteHeader(http.StatusNoContent)
}

func GetChapterByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	chapter, err := repositories.GetChapterByID(context.Background(), id)
	if err != nil {
		http.Error(w, "chapter not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(chapter)
}
