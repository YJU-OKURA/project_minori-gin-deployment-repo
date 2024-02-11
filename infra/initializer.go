package infra

import (
	"github.com/joho/godotenv"
	"log"
)

func Initialize() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file")
	}
}
