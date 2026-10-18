package main

import (
	"carb-backend/cmd/api"
	"carb-backend/cmd/db"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := db.Config{
		Host:     db.Environement.Host,
		Port:     db.Environement.Port,
		User:     db.Environement.User,
		Password: db.Environement.Password,
		Database: db.Environement.Database,
		SSLMode:  db.Environement.SSLMode,
	}

	database, err := db.NewPostgresStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()

	server := api.NewAPIServer(db.Environement.HTTPAddress, database)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM) // Healthier way to end the service
	defer stop()
	if err := server.RunContext(ctx); err != nil {
		log.Fatal(err)
	}
}
