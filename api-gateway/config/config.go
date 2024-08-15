package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	HTTP_PORT            string
	AUTH_SERVICE_PORT    string
	BOOKING_SERVICE_PORT string
	DB_HOST              string
	DB_PORT              int
	DB_USER              string
	DB_PASSWORD          string
	DB_NAME              string
	ACCESS_TOKEN         string
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env: %v", err)
	}
	cfg := Config{}

	cfg.HTTP_PORT = cast.ToString(coalesce("HTTP_PORT", ":8080"))
	cfg.AUTH_SERVICE_PORT = cast.ToString(coalesce("AUTH_SERVICE_PORT", ":50051"))
	cfg.BOOKING_SERVICE_PORT = cast.ToString(coalesce("BOOKING_SERVICE_PORT", ":50052"))

	cfg.DB_HOST = cast.ToString(coalesce("DB_HOST", "localhost"))
	cfg.DB_PORT = cast.ToInt(coalesce("DB_PORT", 5432))
	cfg.DB_USER = cast.ToString(coalesce("DB_USER", "postgres"))
	cfg.DB_PASSWORD = cast.ToString(coalesce("DB_PASSWORD", "123321"))
	cfg.DB_NAME = cast.ToString(coalesce("DB_NAME", "auth_i"))

	cfg.ACCESS_TOKEN = cast.ToString(coalesce("ACCESS_TOKEN", "access_key"))

	return &cfg
}

func coalesce(key string, value interface{}) interface{} {
	val, exists := os.LookupEnv(key)
	if exists {
		return val
	}
	return value
}
