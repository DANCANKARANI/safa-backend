package database

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
    if os.Getenv("GO_ENV") == "test" {
        return nil // Skip database connection in test environment
    }

    _ = godotenv.Load(".env") // Load .env, don't fail if missing

    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        log.Println("Warning: DATABASE_URL is missing. Skipping DB connection.")
        return nil
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Printf("Failed to connect to the PostgreSQL database: %v\n", err)
        return nil
    }

    fmt.Println("PostgreSQL database connection established successfully.")
    return db
}
