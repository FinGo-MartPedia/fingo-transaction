package helpers

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	envMap map[string]string
	once   sync.Once
)

// loadEnv loads .env file into envMap, only once (thread-safe)
func loadEnv() {
	envMap = map[string]string{}

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("[WARNING] .env file not found or failed to load. Using OS environment variables.")
		return
	}

	envMap, err = godotenv.Read(".env")
	if err != nil {
		log.Println("[WARNING] Failed to read .env file content. Using OS environment variables.")
	}
}

// GetEnv fetches environment variable from envMap or OS environment
func GetEnv(key string, defaultValue string) string {
	once.Do(loadEnv)

	// Check from loaded .env map
	if val, ok := envMap[key]; ok && val != "" {
		return val
	}

	// Check from OS environment variable
	if val := os.Getenv(key); val != "" {
		return val
	}

	// Fallback to default value
	return defaultValue
}
