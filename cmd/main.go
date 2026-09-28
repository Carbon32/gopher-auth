package main

import (
	"carb-backend/cmd/api"
	"carb-backend/cmd/db"
	"log"
)

func main() {
	cfg := db.Config{
		Host:     db.Environement.Host,
		Port:     db.Environement.Port,
		User:     db.Environement.User,
		Password: db.Environement.Password,
		Database: db.Environement.Database,
	}

	database, err := db.NewPostgresStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()

	server := api.NewAPIServer(":8080", database)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
