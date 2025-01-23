package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`
	//storage
}

type StorageConfig struct {
}

func InitConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Config path is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Config path does not exist")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error while reading config: %s", err)
	}

	return &cfg
}
