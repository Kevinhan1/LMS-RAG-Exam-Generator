package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"backendLMS/middlewares"
	"backendLMS/models"
	"backendLMS/repositories"

	"github.com/gorilla/mux"
)

// Upload PDF to Supabase Storage via REST API
func uploadPDFToSupabase(file io.Reader, filename string) (string, error) {
	url := fmt.Sprintf("%s/storage/v1/object/materials/%s", os.Getenv("SUPABASE_URL"), filename)
	req, err := http.NewRequest("POST", url, file)
	if err != nil {
		return "", err
	}

	req.Header.Set("apikey", os.Getenv("SUPABASE_SERVICE_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_SERVICE_KEY"))
	req.Header.Set("Content-Type", "application/pdf")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: %s", string(body))
	}

	fileURL := fmt.Sprintf("%s/storage/v1/object/public/materials/%s", os.Getenv("SUPABASE_URL"), filename)
	return fileURL, nil
}

// POST /materials
func CreateMaterial(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)

	// Parse form data
	r.ParseMultipartForm(10 << 20) // max 10MB

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileURL, err := uploadPDFToSupabase(file, header.Filename)
	if err != nil {
		http.Error(w, "failed upload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	material := models.Material{
		TeacherID:   userID,
		CourseID:    int64(atoi(r.FormValue("course_id"))),
		ChapterID:   int64(atoi(r.FormValue("chapter_id"))),
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		FileURL:     fileURL,
		UploadedAt:  time.Now(),
	}

	err = repositories.CreateMaterial(context.Background(), &material)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Log activity
	repositories.CreateLog(context.Background(), &models.LogActivity{
		UserID:      userID,
		Action:      "create_material",
		TargetTable: "materials",
		TargetID:    material.ID,
		Description: material.Title,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(material)
}

// GET /materials
func GetMaterials(w http.ResponseWriter, r *http.Request) {
	materials, _ := repositories.GetMaterials(context.Background())
	json.NewEncoder(w).Encode(materials)
}

func atoi(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func GetMaterialByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	material, err := repositories.GetMaterialByID(context.Background(), id)
	if err != nil {
		http.Error(w, "material not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(material)
}

func UpdateMaterial(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var m models.Material
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	m.ID = id

	if err := repositories.UpdateMaterial(context.Background(), &m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// log
	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	repositories.CreateLog(context.Background(), &models.LogActivity{
		UserID:      userID,
		Action:      "update_material",
		TargetTable: "materials",
		TargetID:    id,
		Description: m.Title,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func DeleteMaterial(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := repositories.DeleteMaterial(context.Background(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	repositories.CreateLog(context.Background(), &models.LogActivity{
		UserID:      userID,
		Action:      "delete_material",
		TargetTable: "materials",
		TargetID:    id,
	})

	w.WriteHeader(http.StatusNoContent)
}
