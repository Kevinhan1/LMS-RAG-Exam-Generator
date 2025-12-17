package main

import (
	"log"
	"net/http"
	"os"

	"backendLMS/db"
	"backendLMS/router"

	"github.com/joho/godotenv"
)

func main() {
	// load env
	_ = godotenv.Load()

	// init database pool
	db.Init()
	defer db.Pool.Close()

	// init router
	r := router.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}