package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env tidak ditemukan, menggunakan environment bawaan sistem")
	}
}

func Get(key string) string {
	return os.Getenv(key)
}
