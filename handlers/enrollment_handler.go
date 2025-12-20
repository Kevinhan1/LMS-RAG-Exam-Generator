package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"backendLMS/middlewares"
	"backendLMS/repositories"

	"github.com/gorilla/mux"
)

func EnrollCourse(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	courseID, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	if err := repositories.Enroll(r.Context(), userID, courseID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func MyCourses(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)

	data, err := repositories.GetUserCourses(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func UnenrollCourse(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.CtxUserID).(int64)
	courseID, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	if err := repositories.Unenroll(r.Context(), userID, courseID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetCourseStudents(w http.ResponseWriter, r *http.Request) {
	courseID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid course id", http.StatusBadRequest)
		return
	}

	data, err := repositories.GetCourseStudents(r.Context(), courseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
