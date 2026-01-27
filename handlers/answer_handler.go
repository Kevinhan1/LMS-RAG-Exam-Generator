package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backendLMS/repositories"

	"github.com/gorilla/mux"
)

func CreateAnswer(w http.ResponseWriter, r *http.Request) {
	questionID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	var req struct {
		Label     string `json:"label"`
		Text      string `json:"text"`
		IsCorrect bool   `json:"is_correct"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err = repositories.CreateAnswer(r.Context(), questionID, repositories.AnswerInput{
		Label:     req.Label,
		Text:      req.Text,
		IsCorrect: req.IsCorrect,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateAnswer(w http.ResponseWriter, r *http.Request) {
	answerID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid id", 400)
		return
	}

	var req struct {
		Label     string `json:"label"`
		Text      string `json:"text"`
		IsCorrect bool   `json:"is_correct"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	err = repositories.UpdateAnswer(r.Context(), answerID, repositories.AnswerInput{
		Label:     req.Label,
		Text:      req.Text,
		IsCorrect: req.IsCorrect,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	answerID, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	err := repositories.DeleteAnswer(r.Context(), answerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
