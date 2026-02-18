package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/imakheri/notifications-thch/config"
	"github.com/imakheri/notifications-thch/internal/domain/entities"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository"
)

func main() {
	cfg := config.Load()
	migrate := flag.Bool("migrate", false, "Create or update database tables with GORM")
	seed := flag.Bool("seed", false, "Seed 'channels' tables with SQL Query")
	all := flag.Bool("all", false, "Executes migration then seeding")

	flag.Parse()

	db := repository.NewDatabase(cfg)

	if *migrate || *all {
		fmt.Println("Executing AutoMigrate...")
		err := db.DatabaseConnection.AutoMigrate(&entities.User{}, &entities.Channel{}, &entities.Notification{})
		if err != nil {
			log.Fatalf("AutoMigrate error: %v", err)
		}
		fmt.Println("Database synchronized")
	}

	if *seed || *all {
		fmt.Println("Seeding 'channels' table...")
		script, err := os.ReadFile("cmd/db/seeds/001_populate_channels.sql")
		if err != nil {
			log.Fatalf("Error reading SQL file: %v", err)
		}

		if err := db.DatabaseConnection.Exec(string(script)).Error; err != nil {
			log.Fatalf("Error seeding table: %v", err)
		}
		fmt.Println("'Channels' table updated")
	}

	if !*migrate && !*seed && !*all {
		fmt.Println("⚠️  To use: go run cmd/db/main.go -migrate | -seed | -all")
	}
}
