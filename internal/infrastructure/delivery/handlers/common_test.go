package handlers

import (
	"os"
	"testing"

	"github.com/imakheri/notifications-thch/config"
	"github.com/imakheri/notifications-thch/internal/infrastructure/repository"
)

var TestDbStruct *repository.Database

func TestMain(m *testing.M) {
	cfg := config.Load()

	TestDbStruct = &repository.Database{
		DatabaseName:     cfg.DatabaseName,
		DatabaseUser:     cfg.DatabaseUser,
		DataBasePassword: cfg.DataBasePassword,
		DatabasePath:     cfg.DatabaseHost,
		DatabasePort:     cfg.DatabasePort,
	}

	TestDbStruct.DatabaseConnection = TestDbStruct.Connect()

	exitCode := m.Run()
	os.Exit(exitCode)
}
