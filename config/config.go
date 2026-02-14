package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseName     string
	DatabaseUser     string
	DataBasePassword string
	DatabasePath     string
	DatabasePort     string
	SecretJWT        string
	ApiKey           string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("could not read env file")
	}

	return &Config{
		DatabaseName:     getEnv("DB_NAME"),
		DatabaseUser:     getEnv("DB_USER"),
		DataBasePassword: getEnv("DB_PASSWORD"),
		DatabasePath:     getEnv("DB_PATH"),
		DatabasePort:     getEnv("DB_PORT"),
		SecretJWT:        getEnv("SECRET_JWT"),
		ApiKey:           getEnv("API_KEY"),
	}
}

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatal("needed environment variable not found")
	}
	return value
}
