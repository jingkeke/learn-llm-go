package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"github.com/smhanov/auth"
)

func main() {
	// Enable WAL mode for concurrency
	db, err := sqlx.Open("sqlite3", "users.db?_journal=WAL&_synchronous=NORMAL")
	if err != nil {
		log.Fatal(err)
	}

	settings := auth.DefaultSettings
	// Note: You would configure SMTP here for real usage.

	// Create the auth handler
	authHandler := auth.New(auth.NewUserDB(db), settings)

	// Configure CORS for Next.js (port 3000) and Node.js backend (port 4000) if they make direct calls
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:4000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
	})

	// Wrap the authHandler with CORS
	handler := c.Handler(authHandler)

	// Ensure that the auth endpoints are available at /user/...
	http.Handle("/user/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Go Auth Service running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
