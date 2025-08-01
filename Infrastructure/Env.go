package infrastructure

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)
type env struct{
    MONGODB_URI string
	DB_NAME string
	Jwt_Secret string
}
var Env env
func LoadEnv() {
	if err := godotenv.Load(); err != nil{
		log.Fatal("could not load the env file")
	}
    Env.DB_NAME = os.Getenv("DB_NAME") 
    Env.Jwt_Secret = os.Getenv("Jwt_Secret") 
    Env.MONGODB_URI = os.Getenv("MONGODB_URI") 
}
