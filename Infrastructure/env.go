package infrastructure

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvStruct struct {
	MONGODB_URI    string
	DB_NAME        string
	JWT_SECRET     string
	PORT           string
	EMAIL_FROM     string
	EMAIL_PORT     string
	EMAIL_HOST     string
	EMAIL_USERNAME string
	EMAIL_PASSWORD string
}

var Env EnvStruct

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env file")
	}

	Env = EnvStruct{
		MONGODB_URI:    os.Getenv("MONGODB_URI"),
		DB_NAME:        os.Getenv("DB_NAME"),
		JWT_SECRET:     os.Getenv("JWT_SECRET"),
		PORT:           os.Getenv("PORT"),
		EMAIL_FROM:     os.Getenv("EMAIL_FROM"),
		EMAIL_HOST:     os.Getenv("EMAIL_HOST"),
		EMAIL_PORT:     os.Getenv("EMAIL_PORT"),
		EMAIL_USERNAME: os.Getenv("EMAIL_USERNAME"),
		EMAIL_PASSWORD: os.Getenv("EMAIL_PASSWORD"),
	}

	if Env.MONGODB_URI == "" || Env.JWT_SECRET == "" || Env.DB_NAME == "" {
		log.Fatal("Missing required environment variables")
	}
}
