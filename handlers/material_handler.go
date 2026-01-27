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
	"bytes"
	"mime/multipart"

	"backendLMS/middlewares"
	"backendLMS/models"
	"backendLMS/repositories"

	"github.com/gorilla/mux"
)

/*
====================================
 Supabase Storage Upload (PUT)
====================================
*/
func uploadPDFToSupabase(file io.Reader, filename string) (string, error) {
	url := fmt.Sprintf(
		"%s/storage/v1/object/materials/%s",
		os.Getenv("SUPABASE_URL"),
		filename,
	)

	req, err := http.NewRequest(http.MethodPut, url, file)
	if err != nil {
		return "", err
	}

	req.Header.Set("apikey", os.Getenv("SUPABASE_SERVICE_KEY"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_SERVICE_KEY"))
	req.Header.Set("Content-Type", "application/pdf")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("supabase upload failed: %s", string(body))
	}

	// public read URL (TANPA service key)
	fileURL := fmt.Sprintf(
		"%s/storage/v1/object/public/materials/%s",
		os.Getenv("SUPABASE_URL"),
		filename,
	)

	return fileURL, nil
}

/*
====================================
 POST /materials
====================================
*/
func CreateMaterial(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid multipart form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	/*
	--------------------------------
	 Validasi file type (WAJIB)
	--------------------------------
	*/
	if header.Header.Get("Content-Type") != "application/pdf" {
		http.Error(w, "only PDF allowed", http.StatusBadRequest)
		return
	}

	// Magic number check (%PDF)
	buf := make([]byte, 4)
	if _, err := file.Read(buf); err != nil || string(buf) != "%PDF" {
		http.Error(w, "invalid PDF file", http.StatusBadRequest)
		return
	}

	// reset reader
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}

	/*
	--------------------------------
	 Anti filename collision
	--------------------------------
	*/
	filename := fmt.Sprintf(
		"%d_%d_%s",
		userID,
		time.Now().Unix(),
		header.Filename,
	)

	fileURL, err := uploadPDFToSupabase(file, filename)
	if err != nil {
		http.Error(w, "upload failed: "+err.Error(), http.StatusInternalServerError)
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

	if err := repositories.CreateMaterial(context.Background(), &material); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ⬇️ kirim PDF ke FastAPI untuk chunking + embedding
	if _, err := file.Seek(0, io.SeekStart); err == nil {
		go sendPDFToFastAPI(
			file,
			filename,
			material.ID,
			material.CourseID,
			material.ChapterID,
		)
	}

	// log activity
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

/*
====================================
 GET /materials
====================================
*/
func GetMaterials(w http.ResponseWriter, r *http.Request) {
	materials, err := repositories.GetMaterials(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(materials)
}

/*
====================================
 GET /materials/{id}
====================================
*/
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

	json.NewEncoder(w).Encode(material)
}

/*
====================================
 PUT /materials/{id}
====================================
*/
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

	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	repositories.CreateLog(context.Background(), &models.LogActivity{
		UserID:      userID,
		Action:      "update_material",
		TargetTable: "materials",
		TargetID:    id,
		Description: m.Title,
	})

	json.NewEncoder(w).Encode(m)
}

/*
====================================
 PUT /teacher/materials/{id}
====================================
*/
func TeacherUpdateMaterial(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)

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

	if err := repositories.UpdateMaterialByTeacher(
		context.Background(),
		&m,
		userID,
	); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	repositories.CreateLog(context.Background(), &models.LogActivity{
		UserID:      userID,
		Action:      "update_material",
		TargetTable: "materials",
		TargetID:    id,
		Description: m.Title,
	})

	json.NewEncoder(w).Encode(m)
}

/*
====================================
 DELETE /materials/{id}
====================================
*/
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

/*
====================================
 DELETE /teacher/materials/{id}
====================================
*/
func TeacherDeleteMaterial(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)

	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := repositories.DeleteMaterialByTeacher(
		context.Background(),
		id,
		userID,
	); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	repositories.CreateLog(context.Background(), &models.LogActivity{
		UserID:      userID,
		Action:      "delete_material",
		TargetTable: "materials",
		TargetID:    id,
	})

	w.WriteHeader(http.StatusNoContent)
}

func sendPDFToFastAPI(
	file io.Reader,
	filename string,
	materialID int64,
	courseID int64,
	chapterID int64,
) error {

	url := os.Getenv("FASTAPI_URL") + "/ingest_material"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// file
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, file); err != nil {
		return err
	}

	// metadata
	writer.WriteField("material_id", fmt.Sprintf("%d", materialID))
	writer.WriteField("course_id", fmt.Sprintf("%d", courseID))
	writer.WriteField("chapter_id", fmt.Sprintf("%d", chapterID))

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("fastapi ingest failed: %s", string(b))
	}

	return nil
}

/*
====================================
 Utils
====================================
*/
func atoi(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
