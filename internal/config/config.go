package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	ServerAddress string
	LogLevel      string
	AdminEmail    string
	AdminPassword string
}

// LoadConfig loads configuration from environment variables or .env file
func LoadConfig() *Config {
	viper.AutomaticEnv()

	// Attempt to read configuration from environment variables
	//cfg := readFromEnv()

	// If configuration is not found in environment variables, attempt to read from .env file
	cfg := readFromEnvFile()

	return cfg
}

func readFromEnv() *Config {
	return &Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		LogLevel:      os.Getenv("LOG_LEVEL"),
		AdminEmail:    os.Getenv("ADMIN_EMAIL"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
	}
}

func readFromEnvFile() *Config {
	// Check if .env file exists
	if _, err := os.Stat(".env"); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
	}

	// Read configuration from .env file using viper
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		return nil
	}

	// Initialize config struct
	cfg := &Config{
		DBHost:        viper.GetString("DB_HOST"),
		DBPort:        viper.GetString("DB_PORT"),
		DBUser:        viper.GetString("DB_USER"),
		DBPassword:    viper.GetString("DB_PASSWORD"),
		DBName:        viper.GetString("DB_NAME"),
		JWTSecret:     viper.GetString("JWT_SECRET"),
		ServerAddress: viper.GetString("SERVER_ADDRESS"),
		LogLevel:      viper.GetString("LOG_LEVEL"),
		AdminEmail:    viper.GetString("ADMIN_EMAIL"),
		AdminPassword: viper.GetString("ADMIN_PASSWORD"),
	}

	return cfg
}
