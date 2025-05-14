package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"zero-agency-tambov/config"
	"zero-agency-tambov/logger"
)

func ConnectDB(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Log.Fatalf("Error connecting to DB: %v", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	if err := db.Ping(); err != nil {
		logger.Log.Fatalf("Failed to check connection to database: %v", err)
	}

	logger.Log.Info("Successful connection to the database!")
	return db
}
