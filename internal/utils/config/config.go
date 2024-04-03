package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	ServerPort string
	DbPassword string
	DbName     string
	DbUser     string
	Secret     string
}

func NewConfig() *Config {
	return &Config{}
}

func (cfg *Config) InitENV() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg.ServerPort = os.Getenv("SERVER_PORT")
	cfg.DbPassword = os.Getenv("DB_PASSWORD")
	cfg.DbName = os.Getenv("DB_NAME")
	cfg.DbUser = os.Getenv("DB_USER")
	cfg.Secret = os.Getenv("SECRET_FOR_TOKEN")

}
