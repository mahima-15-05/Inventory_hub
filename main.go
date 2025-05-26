package main

import (
	"inventory-control-hub/database"
	"inventory-control-hub/routes"

	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file...")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080" // default port
	}
	database.Connect()
	database.Migrate()
	log.Println("Server running at http://localhost:" + port)
	

	http.ListenAndServe(":"+port, routes.SetupRouter())
}
