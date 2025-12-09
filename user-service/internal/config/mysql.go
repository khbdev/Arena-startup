package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// NewMySQLConnection - .env orqali o'qib, DB ga ulanish
func NewMySQLConnection() *sql.DB {
	// env fayldan o'qish
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbPort == "" {
		dbPort = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to open MySQL connection: %v", err)
	}

	// connection test
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping MySQL: %v", err)
	}

	return db
}
