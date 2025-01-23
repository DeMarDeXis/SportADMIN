package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	JsonPath string `mapstructure:"jsonPath"`
	DB       DB     `mapstructure:"db"`
}

type DB struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	DBName   string `mapstructure:"db_name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

func LoadConfig(defaultPath string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	configPath := os.Getenv("CONFIG_PATH")
	println(configPath)
	if configPath == "" {
		configPath = defaultPath
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	var config Config

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
