package router

import (
	"net/http"

	"backendLMS/handlers"
	"backendLMS/middlewares"

	"github.com/gorilla/mux"
)

func New() http.Handler {
	r := mux.NewRouter()

	// ======================
	// PUBLIC ROUTES
	// ======================
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods("GET")

	// ======================
	// PROTECTED ROUTES (JWT)
	// ======================
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middlewares.JWTAuth)

	// ======================
	// COMMON USER
	// ======================
	api.HandleFunc("/me", handlers.Me).Methods("GET")

	// ======================
	// ADMIN ONLY
	// ======================
	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(middlewares.RequireRoles(1)) // ADMIN ONLY

	// ---- Role Management
	admin.HandleFunc("/roles", handlers.GetRoles).Methods("GET")
	admin.HandleFunc("/roles", handlers.CreateRole).Methods("POST")
	admin.HandleFunc("/roles/{id}", handlers.GetRole).Methods("GET")
	admin.HandleFunc("/roles/{id}", handlers.UpdateRole).Methods("PUT")
	admin.HandleFunc("/roles/{id}", handlers.DeleteRole).Methods("DELETE")

	// ---- User Management
	admin.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	admin.HandleFunc("/users/{id}", handlers.UpdateUserRole).Methods("PUT")
	admin.HandleFunc("/users/{id}", handlers.DeleteUser).Methods("DELETE")

	// ---- Register
	admin.HandleFunc("/register/student", handlers.RegisterStudent).Methods("POST")
	admin.HandleFunc("/register/teacher", handlers.RegisterTeacher).Methods("POST")

	// ---- Course Management
	admin.HandleFunc("/courses", handlers.CreateCourse).Methods("POST")
	admin.HandleFunc("/courses/{id}", handlers.UpdateCourse).Methods("PUT")
	admin.HandleFunc("/courses/{id}", handlers.DeleteCourse).Methods("DELETE")

	// ======================
	// PUBLIC / AUTH USERS
	// ======================
	api.HandleFunc("/courses", handlers.GetCourses).Methods("GET")
	api.HandleFunc("/courses/{id}", handlers.GetCourseByID).Methods("GET")

	// ---- Class
	admin.HandleFunc("/classes", handlers.CreateClass).Methods("POST")
	admin.HandleFunc("/classes/{id}", handlers.UpdateClass).Methods("PUT")
	admin.HandleFunc("/classes/{id}", handlers.DeleteClass).Methods("DELETE")
	api.HandleFunc("/classes", handlers.GetClasses).Methods("GET")
	api.HandleFunc("/classes/{id}", handlers.GetClassByID).Methods("GET")

	// ---- Chapter
	admin.HandleFunc("/chapters", handlers.CreateChapterHandler).Methods("POST")
	admin.HandleFunc("/chapters/{id}", handlers.UpdateChapterHandler).Methods("PUT")
	admin.HandleFunc("/chapters/{id}", handlers.DeleteChapterHandler).Methods("DELETE")
	api.HandleFunc("/courses/{course_id}/chapters", handlers.GetChaptersHandler).Methods("GET")
	api.HandleFunc("/chapters/{id}", handlers.GetChapterByIDHandler).Methods("GET")

	// ---- Tag
	admin.HandleFunc("/tags", handlers.CreateTagHandler).Methods("POST")
	admin.HandleFunc("/tags/{id}", handlers.UpdateTagHandler).Methods("PUT")
	admin.HandleFunc("/tags/{id}", handlers.DeleteTagHandler).Methods("DELETE")
	api.HandleFunc("/tags", handlers.GetTagsHandler).Methods("GET")
	api.HandleFunc("/tags/{id}", handlers.GetTagByID).Methods("GET")

	// ---- Material
	admin.HandleFunc("/materials", handlers.CreateMaterial).Methods("POST")
	admin.HandleFunc("/materials/{id}", handlers.UpdateMaterial).Methods("PUT")
	admin.HandleFunc("/materials/{id}", handlers.DeleteMaterial).Methods("DELETE")
	api.HandleFunc("/materials", handlers.GetMaterials).Methods("GET")
	api.HandleFunc("/materials/{id}", handlers.GetMaterialByID).Methods("GET")

	// ======================
	// ENROLLMENT (AUTH FIXED)
	// ======================

	// STUDENT ONLY
	api.Handle(
		"/courses/{id}/enroll",
		middlewares.RequireRoles(3)(
			http.HandlerFunc(handlers.EnrollCourse),
		),
	).Methods("POST")

	api.Handle(
		"/courses/{id}/enroll",
		middlewares.RequireRoles(3)(
			http.HandlerFunc(handlers.UnenrollCourse),
		),
	).Methods("DELETE")

	api.Handle(
		"/my-courses",
		middlewares.RequireRoles(3)(
			http.HandlerFunc(handlers.MyCourses),
		),
	).Methods("GET")

	// ======================
	// COURSE STUDENTS
	// ======================

	// TEACHER + ADMIN
	api.Handle(
		"/courses/{id}/students",
		middlewares.RequireRoles(2, 1)(
			http.HandlerFunc(handlers.GetCourseStudents),
		),
	).Methods("GET")

	return r
}