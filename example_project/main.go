package main

import (
	"log"
	"os"

	"github.com/ralvescosta/dotenv"
)

func main() {
	dotenv.Configure(".env.development")

	log.Printf("DB HOST: %v \n", os.Getenv("DB_HOST"))
	log.Printf("DB USER: %v \n", os.Getenv("DB_USER"))
	log.Printf("DB PASS: %v \n", os.Getenv("DB_PASS"))
}
