package main

import (
	"log"
	"os"
	"something/config"
	"something/database"
	"something/server"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DSN")

	db, err := storage.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	panic(server.NewServer(
		db,
		&config.AuthConfig{
			AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
			RefreshTokenSecret: os.Getenv("REFRESH_TOKEN_SECRET"),
			AccessTokenTTL: 60 * time.Minute,
			RefreshTokenTTL: 24 * time.Hour,
		},
	).Run())
}