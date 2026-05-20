package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/impulseweb3/hypixel-skyblock-auctions-go/internal/hypixel"
	"github.com/impulseweb3/hypixel-skyblock-auctions-go/internal/persistent"
	"github.com/impulseweb3/hypixel-skyblock-auctions-go/internal/worker"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		slog.Error(err.Error())
		return
	}

	if os.Getenv("ENVIRONMENT") == "production" {
		slog.SetLogLoggerLevel(slog.LevelInfo)

	} else {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DATABASE")
	port := os.Getenv("POSTGRES_PORT")
	sslmode := os.Getenv("POSTGRES_SSLMODE")
	timezone := os.Getenv("POSTGRES_TIMEZONE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone,
	)

	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		slog.Error(err.Error())
		return
	}

	database := persistent.NewDatabase(db)

	err = database.AutoMigrate()

	if err != nil {
		slog.Error(err.Error())
		return
	}

	repository := persistent.NewRepository(database)
	hypixelClient := hypixel.NewClient(&http.Client{})

	tracker := worker.NewTracker(repository, hypixelClient)
	tracker.Start()
}
