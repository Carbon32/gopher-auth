package main

import (
	"carb-backend/cmd/db"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://cmd/migrate/migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Usage: migrate (up/down)")
	}

	cmd := os.Args[1]

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		} else {
			log.Println("Success")
		}

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		} else {
			log.Println("Success")
		}

	default:
		log.Fatalf("Unknown command: %s", cmd)
	}
}
