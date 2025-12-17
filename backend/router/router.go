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
	r.HandleFunc("/register", handlers.Register).Methods("POST")
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
	// ADMIN ONLY
	// ======================
	api.Handle("/roles",
		middlewares.RequireRole(1)(
			http.HandlerFunc(handlers.GetRoles),
		),
	).Methods("GET")

	api.Handle("/roles",
		middlewares.RequireRole(1)(
			http.HandlerFunc(handlers.CreateRole),
		),
	).Methods("POST")

	return r
}
