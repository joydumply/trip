package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	RedisAddr   string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, trying environment variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, ErrMissingEnv("DATABASE_URL")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		return nil, ErrMissingEnv("REDIS_ADDR")
	}

	return &Config{
		DatabaseURL: dbURL,
		RedisAddr:   redisAddr,
	}, nil
}

type ErrMissingEnv string

func (e ErrMissingEnv) Error() string {
	return "missing env variable: " + string(e)
}
