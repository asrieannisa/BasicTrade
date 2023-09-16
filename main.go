package main

import (
	"BasicTrade/database"
	router "BasicTrade/routers"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load variables from .env file
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	// Get the PORT value from the environment
	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "7070" // Default port if not specified in .env
	}

	database.StartDB()
	r := router.StartApp()
	r.Run(":" + PORT)
}
