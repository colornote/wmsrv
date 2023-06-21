package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// test env
func ExampleEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Printf("%++v\n", os.Getenv("ENABLE_CSRF"))

	// Output: test env
}
