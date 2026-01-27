package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"backendLMS/middlewares"
	"backendLMS/repositories"
)

var FASTAPI_URL = os.Getenv("FASTAPI_URL")

type createQuestionRequest struct {
	MaterialID    int64  `json:"material_id"`
	Content       string `json:"content"`
	Difficulty    string `json:"difficulty"`
	TaxonomyLevel string `json:"taxonomy_level"`
	Answers       []struct {
		Label     string `json:"label"`
		Text      string `json:"text"`
		IsCorrect bool   `json:"is_correct"`
	} `json:"answers"`
}

type fastAPIGeneratedQuestion struct {
	MaterialID    int64  `json:"material_id"`
	Content       string `json:"content"`
	Difficulty    string `json:"difficulty"`
	TaxonomyLevel string `json:"taxonomy_level"`
	Answers       []struct {
		Label     string `json:"label"`
		Text      string `json:"text"`
		IsCorrect bool   `json:"is_correct"`
	} `json:"answers"`
}

type RAGGenerateRequest struct {
	MaterialID  int64  `json:"material_id"`
	Instruction string `json:"instruction"`
}

func GenerateQuestionFromRAG(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)

	var req struct {
		MaterialID  int64  `json:"material_id"`
		Instruction string `json:"instruction"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// 1. Panggil FastAPI
	payload := RAGGenerateRequest{
		MaterialID:  req.MaterialID,
		Instruction: req.Instruction,
	}
	body, _ := json.Marshal(payload)

	httpReq, _ := http.NewRequest("POST", FASTAPI_URL+"/generate_exam", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "failed to generate question from RAG", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var fastApiResp struct {
		RagResult []struct {
			Content       string `json:"content"`
			Difficulty    string `json:"difficulty"`
			TaxonomyLevel string `json:"taxonomy_level"`
			Answers       []struct {
				Label     string `json:"label"`
				Text      string `json:"text"`
				IsCorrect bool   `json:"is_correct"`
			} `json:"answers"`
		} `json:"rag_result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&fastApiResp); err != nil {
		http.Error(w, "invalid response from RAG", http.StatusInternalServerError)
		return
	}

	// 2. Simpan ke DB (Batch Insert)
	for _, q := range fastApiResp.RagResult {
		var answers []repositories.AnswerInput
		for _, a := range q.Answers {
			answers = append(answers, repositories.AnswerInput{
				Label:     a.Label,
				Text:      a.Text,
				IsCorrect: a.IsCorrect,
			})
		}

		err = repositories.CreateQuestionWithAnswers(
			r.Context(),
			req.MaterialID,
			userID,
			q.Content,
			q.Difficulty,
			q.TaxonomyLevel,
			answers,
		)
		if err != nil {
			// Jika satu gagal, return error (bisa diimprove dengan bulk insert transaction)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func CreateQuestion(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)

	var req createQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if len(req.Answers) < 2 {
		http.Error(w, "minimum 2 answers", http.StatusBadRequest)
		return
	}

	var correct int
	var answers []repositories.AnswerInput
	for _, a := range req.Answers {
		if a.IsCorrect {
			correct++
		}
		answers = append(answers, repositories.AnswerInput{
			Label:     a.Label,
			Text:      a.Text,
			IsCorrect: a.IsCorrect,
		})
	}

	if correct != 1 {
		http.Error(w, "exactly 1 correct answer required", http.StatusBadRequest)
		return
	}

	err := repositories.CreateQuestionWithAnswers(
		r.Context(),
		req.MaterialID,
		userID,
		req.Content,
		req.Difficulty,
		req.TaxonomyLevel,
		answers,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetQuestions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	roleID := r.Context().Value(middlewares.CtxRoleID).(int64)

	data, err := repositories.GetQuestions(r.Context(), userID, roleID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(data)
}

func GetQuestionDetail(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	roleID := r.Context().Value(middlewares.CtxRoleID).(int64)
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	q, answers, err := repositories.GetQuestionByID(r.Context(), id, userID, roleID)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"question": q,
		"answers":  answers,
	})
}

func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	roleID := r.Context().Value(middlewares.CtxRoleID).(int64)
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	var req createQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	var correct int
	var answers []repositories.AnswerInput
	for _, a := range req.Answers {
		if a.IsCorrect {
			correct++
		}
		answers = append(answers, repositories.AnswerInput{
			Label:     a.Label,
			Text:      a.Text,
			IsCorrect: a.IsCorrect,
		})
	}

	if correct != 1 {
		http.Error(w, "exactly one correct answer required", 400)
		return
	}

	err := repositories.UpdateQuestion(
		r.Context(), id, userID, roleID,
		req.Content, req.Difficulty, req.TaxonomyLevel,
		answers,
	)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	roleID := r.Context().Value(middlewares.CtxRoleID).(int64)
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	err := repositories.DeleteQuestion(r.Context(), id, userID, roleID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateQuestionStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	roleID := r.Context().Value(middlewares.CtxRoleID).(int64)
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	var body struct {
		Status string `json:"status"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	err := repositories.UpdateQuestionStatus(
		r.Context(), id, userID, roleID, body.Status,
	)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusOK)
}
