package utils

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	DatabaseURL     string
	TestDatabaseURL string
	JwtSecret       string
	SecretToken     string
}

func NewConfig(path string) Config {
	projectRoot, err := GetProjectRoot()
	if err != nil {
		panic(fmt.Sprintf("Error finding project root: %v", err))
	}

	envPath := filepath.Join(projectRoot, ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		panic("Error loading .env file")
	}
	return Config{
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		TestDatabaseURL: os.Getenv("TEST_DATABASE_URL"),
		JwtSecret:       os.Getenv("JWT_SECRET"),
		SecretToken:     os.Getenv("SECRET_TOKEN"),
	}
}

func GetEnv(key string) string {
	projectRoot, err := GetProjectRoot()
	if err != nil {
		panic(fmt.Sprintf("Error finding project root: %v", err))
	}

	envPath := filepath.Join(projectRoot, ".env")
	if err := godotenv.Load(envPath); err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	return os.Getenv(key)
}
func GetProjectRoot() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}
	currentDir := filepath.Dir(filename)

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return ".", nil
		}
		currentDir = parentDir
	}
}
