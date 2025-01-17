package main

import (
	"inventoryapp/database"
	"inventoryapp/router"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load("config/.env")

	if err != nil {
		log.Fatal("No file .env found, using environment system variables")
	}

	PORT := os.Getenv("API_PORT")

	if PORT == "" {
		PORT = "8080"
	}

	database.StartDB()
	router.StartServer().Run(":" + PORT)
}
