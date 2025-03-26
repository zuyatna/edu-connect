package database

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var (
	instance    *sql.DB
	oncePostgre sync.Once
)

func GetDB() *sql.DB {
	oncePostgre.Do(func() {
		connectionString := getConnectionString()
		db, err := sql.Open("pgx", connectionString)
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to database: %v", err))
		}

		err = db.Ping()
		if err != nil {
			panic(fmt.Sprintf("Failed to ping database: %v", err))
		}

		instance = db
		fmt.Println("Database connection established successfully")
	})

	return instance
}

func getConnectionString() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: error loading .env file")
	}

	return os.Getenv("POSTGRES_URI")
}

func CloseDB() {
	if instance != nil {
		instance.Close()
		instance = nil
		fmt.Println("Database connection closed")
	}
}
