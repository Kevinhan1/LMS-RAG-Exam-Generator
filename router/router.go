package router

import (
	"net/http"

	"backendLMS/handlers"
	"backendLMS/middlewares"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

	// ---- Question Management (ADMIN FULL CRUD)
	admin.HandleFunc("/questions", handlers.GetQuestions).Methods("GET")
	admin.HandleFunc("/questions/{id}", handlers.GetQuestionDetail).Methods("GET")
	admin.HandleFunc("/questions/{id}", handlers.UpdateQuestion).Methods("PUT")
	admin.HandleFunc("/questions/{id}", handlers.DeleteQuestion).Methods("DELETE")
	admin.HandleFunc("/questions/{id}/status", handlers.UpdateQuestionStatus).Methods("PATCH")

	// ---- Answer Management (ADMIN)
	admin.HandleFunc(
		"/questions/{id}/answers",
		handlers.CreateAnswer,
	).Methods("POST")

	admin.HandleFunc(
		"/answers/{id}",
		handlers.UpdateAnswer,
	).Methods("PUT")

	admin.HandleFunc(
		"/answers/{id}",
		handlers.DeleteAnswer,
	).Methods("DELETE")

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

	// ---- Permission Management

	// READ (AUTH USER)
	api.HandleFunc("/permissions", handlers.GetPermissions).Methods("GET")
	api.HandleFunc("/permissions/{id}", handlers.GetPermissionByID).Methods("GET")

	// WRITE (ADMIN)
	admin.HandleFunc("/permissions", handlers.CreatePermission).Methods("POST")
	admin.HandleFunc("/permissions/{id}", handlers.UpdatePermission).Methods("PUT")
	admin.HandleFunc("/permissions/{id}", handlers.DeletePermission).Methods("DELETE")

	// ---- Role Permission Management
	admin.HandleFunc(
		"/roles/{id}/permissions",
		handlers.AssignPermission,
	).Methods("POST")

	admin.HandleFunc(
		"/roles/{id}/permissions",
		handlers.GetRolePermissions,
	).Methods("GET")

	admin.HandleFunc(
		"/roles/{id}/permissions/{permission_id}",
		handlers.RemovePermission,
	).Methods("DELETE")

	// ======================
	// TEACHER ONLY
	// ======================
	teacher := api.PathPrefix("/teacher").Subrouter()
	teacher.Use(middlewares.RequireRoles(2)) // ROLE_TEACHER

	// ---- Material (TEACHER - OWN ONLY)
	teacher.HandleFunc(
		"/materials/{id}",
		handlers.TeacherUpdateMaterial,
	).Methods("PUT")

	teacher.HandleFunc(
		"/materials/{id}",
		handlers.TeacherDeleteMaterial,
	).Methods("DELETE")

	teacher.HandleFunc(
		"/materials",
		handlers.CreateMaterial,
	).Methods("POST")

	teacher.HandleFunc("/questions", handlers.GetQuestions).Methods("GET")
	teacher.HandleFunc("/questions/{id}", handlers.GetQuestionDetail).Methods("GET")
	teacher.HandleFunc("/questions/{id}", handlers.UpdateQuestion).Methods("PUT")
	teacher.HandleFunc("/questions/{id}", handlers.DeleteQuestion).Methods("DELETE")
	teacher.HandleFunc("/questions/{id}/status", handlers.UpdateQuestionStatus).Methods("PATCH")

	// ---- Answer Management (TEACHER)
	teacher.HandleFunc(
		"/questions/{id}/answers",
		handlers.CreateAnswer,
	).Methods("POST")

	teacher.HandleFunc(
		"/answers/{id}",
		handlers.UpdateAnswer,
	).Methods("PUT")

	teacher.HandleFunc(
		"/answers/{id}",
		handlers.DeleteAnswer,
	).Methods("DELETE")

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
	// ---- Material (ADMIN)
	admin.HandleFunc("/materials/{id}", handlers.UpdateMaterial).Methods("PUT")
	admin.HandleFunc("/materials/{id}", handlers.DeleteMaterial).Methods("DELETE")
	api.HandleFunc("/materials", handlers.GetMaterials).Methods("GET")
	api.HandleFunc("/materials/{id}", handlers.GetMaterialByID).Methods("GET")
	admin.HandleFunc(
		"/materials",
		handlers.CreateMaterial,
	).Methods("POST")

	// ======================
	// ENROLLMENT (AUTH FIXED
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

	// ---- Material Tags
	admin.HandleFunc(
		"/materials/{id}/tags",
		handlers.AttachTagToMaterial,
	).Methods("POST")

	admin.HandleFunc(
		"/materials/{id}/tags/{tag_id}",
		handlers.DetachMaterialTag,
	).Methods("DELETE")

	api.HandleFunc(
		"/materials/{id}/tags",
		handlers.GetMaterialTags,
	).Methods("GET")

	teacher.HandleFunc(
		"/questions/rag_generate",
		handlers.GenerateQuestionFromRAG,
	).Methods("POST")

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Change this to specific domain in production
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	return c.Handler(r)
}
