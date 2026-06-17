package database

import (
	"crescendo-api/config/env"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

func NewConnection() (*sql.DB, error) {
	var err error
	err = env.Load()
	if err != nil {
		log.Fatalf("Unable to load environmental variables file: %v", err)
	}
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	sslMode := os.Getenv("POSTGRES_SSL_MODE")
	sslRootCert := os.Getenv("POSTGRES_SSL_ROOT_CERT")

	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbname,
		sslMode,
	)

	if sslRootCert != "" {
		dsn += fmt.Sprintf(" sslrootcert=%s", sslRootCert)
	}

	DB, err := sql.Open("postgres", dsn)

	// Config del pool de conexiones
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxIdleTime(15 * time.Minute)
	DB.SetConnMaxLifetime(time.Hour)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}
	return DB, nil
}
