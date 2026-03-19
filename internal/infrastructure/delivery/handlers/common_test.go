package handlers

import (
	"log"
	"os"
	"testing"

	"github.com/imakheri/notifications-thch/internal/infrastructure/repository"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository/dtos"
)

var TestDbStruct *repository.Database

func TestMain(m *testing.M) {
	TestDbStruct = &repository.Database{
		DatabaseName:     "notification_manager_db_test",
		DatabaseUser:     "root",
		DataBasePassword: "root_password",
		DatabasePath:     "127.0.0.1",
		DatabasePort:     "3307",
	}

	TestDbStruct.DatabaseConnection = TestDbStruct.Connect()

	err := TestDbStruct.DatabaseConnection.AutoMigrate(&dtos.UserModel{}, &dtos.ChannelModel{}, &dtos.NotificationModel{})
	if err != nil {
		log.Fatalf("AutoMigrate error: %v", err)
	}

	script, err := os.ReadFile("../../../../cmd/db/seeds/001_populate_channels.sql")
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}

	if err = TestDbStruct.DatabaseConnection.Exec(string(script)).Error; err != nil {
		log.Fatalf("Error seeding table: %v", err)
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}
