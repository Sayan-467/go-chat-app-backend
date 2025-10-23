package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	DbHost string 
	DbPort string
	DbUser string
	DbPassword string
	DbName string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env found")
	}

	return Config{
		Port: getEnv("PORT", "8080"),
		DbHost: getEnv("DB_HOST", "localhost"),
		DbPort: getEnv("DB_PORT", "5432"),
		DbUser: getEnv("DB_USER", "postgres"),
		DbPassword: getEnv("DB_PASSWORD", "Sayan@postgresql25"),
		DbName: getEnv("DB_NAME", "chatappDB"),
	}
}

func getEnv(key, fallback string) string {
	if val, isExist := os.LookupEnv(key); isExist {
		return  val
	}
	return fallback
}