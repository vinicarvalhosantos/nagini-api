package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func Config(key, fallback string) string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Print("Error loading .env file")
	}

	return getEnv(key, fallback)
}

func getEnv(key, fallback string) string {
	variable := os.Getenv(key)
	if len(variable) == 0 {
		return fallback
	}

	return variable
}

func GetSecretKey(key string) string {
	return Config(key, "stage")
}
