package config

import (
    "log"
    "os"
    "path/filepath"

    "github.com/joho/godotenv"
)

func LoadEnv() {
    // Get the current working directory
    rootDir, err := os.Getwd()
    if err != nil {
        log.Fatalf("Error getting current working directory: %v", err)
    }

    // Construct the path to the environment file
    envPath := filepath.Join(rootDir, "internal", "config", "config.env") 

    // Load the environment file
    err = godotenv.Load(envPath)
    if err != nil {
        log.Fatalf("Error loading config.env file: %v", err)
    }
    log.Println("Env Loaded :)")
}

func GetEnv(key string) string {
    return os.Getenv(key)
}
