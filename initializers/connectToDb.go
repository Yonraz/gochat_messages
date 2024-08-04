package initializers

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB


func ConnectToDb() {
	host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
	fmt.Println(dsn)

	maxRetries := 3
    initialDelay := time.Second

    for i := 0; i < maxRetries; i++ {
        DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            fmt.Println("Connected to postgres!")
            return
        }

        // Log the error and wait before retrying
        fmt.Printf("Failed to connect to postgres: %v. Retrying in %v...\n", err, initialDelay)
        time.Sleep(initialDelay)

        // Exponential backoff
        initialDelay *= 2
    }
	if err != nil {
		panic("Error connecting to postgres")
	}
}