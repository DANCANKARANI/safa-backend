package database

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
    if os.Getenv("GO_ENV") == "test" {
        return nil // Skip database connection in test environment
    }

    // Try to load .env file but don't fail if it's not found
    _ = godotenv.Load(".env")
    
    dbUser := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    
    if dbUser == "" || dbName == "" || password == "" || host == "" || port == "" {
        log.Println("Warning: Database configuration variables are missing. Using test mode.")
        return nil
    }
   
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        dbUser, password, host, port, dbName)
    
    db, err := gorm.Open(mysql.Open(dsn))
    if err != nil {
        log.Printf("Failed to connect to the database: %v\n", err)
        return nil
    }
    
    fmt.Println("Database connection established successfully.")
    return db
}