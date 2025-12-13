package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvFile() string {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "development"
	}

	envFile := fmt.Sprintf(".env.%s", appEnv)
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return ".env"
	}
	return envFile
}

func LoadEnvFile(filename string) error {
	return godotenv.Load(filename)
}
