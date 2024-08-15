package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

type Config struct {
	BOOKING_SERVICE_PORTD string
	BOOKING_SERVICE_PORT  string
	AUTH_SERVICE_PORT     string
	MongoDB_NAME          string
	MongoURI              string
	MongoURIL             string
	RABBIT                string
	RABBITL               string
}

func coalesce(env string, defaultValue interface{}) interface{} {
	value, exists := os.LookupEnv(env)
	if exists {
		return value
	}
	return defaultValue
}

func Load() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := Config{}

	cfg.AUTH_SERVICE_PORT = cast.ToString(coalesce("AUTH_SERVICE_PORT", ":50051"))
	cfg.BOOKING_SERVICE_PORT = cast.ToString(coalesce("BOOKING_SERVICE_PORT", ":50052"))
	cfg.BOOKING_SERVICE_PORTD = cast.ToString(coalesce("BOOKING_SERVICE_PORTD", "booking-service"+":50052"))

	cfg.MongoDB_NAME = cast.ToString(coalesce("MongoDB_NAME", "booking"))
	cfg.MongoURI = cast.ToString(coalesce("MongoURI", "mongodb://localhost:27017"))
	cfg.MongoURIL = cast.ToString(coalesce("MongoURIL", "mongodb://localhost:27017"))

	cfg.RABBIT = cast.ToString(coalesce("RABBIT", "amqp://guest:guest@localhost:5672/"))

	cfg.RABBITL = cast.ToString(coalesce("RABBITL", "amqp://guest:guest@localhost:5672/"))

	return &cfg
}
