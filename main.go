package main

import (
	"apisrv/api"
	"apisrv/api/auth"
	"apisrv/database"
	"apisrv/pkg"
	"apisrv/server"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = pkg.Setup()

	if err != nil {
		log.Fatal(err)
	}

	// Server initialization
	app := server.Create()

	// Migrations
	database.DB.AutoMigrate(
		&auth.User{},
	)

	// Api routes
	api.Setup(app)

	if err := server.Listen(app); err != nil {
		log.Panic(err)
	}
}
