package repository

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/imakheri/notifications-thch/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	DatabaseName     string
	DatabaseUser     string
	DataBasePassword string
	DatabasePath     string
	DatabasePort     string
}

func NewDatabase(cfg *config.Config) *Database {
	return &Database{
		DatabaseName:     cfg.DatabaseName,
		DatabaseUser:     cfg.DatabaseUser,
		DataBasePassword: cfg.DataBasePassword,
		DatabasePath:     cfg.DatabasePath,
		DatabasePort:     cfg.DatabasePort,
	}
}

func (db *Database) Database() *gorm.DB {
	var dns = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v%v",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PATH"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		"?charset=utf8mb4&parseTime=True&loc=Local",
	)
	if db, err := gorm.Open(mysql.Open(dns), &gorm.Config{}); err != nil {
		panic(err)
	} else {
		poolConnection, err := db.DB()
		if err != nil {
			log.Fatal(err)
		}
		poolConnection.SetMaxIdleConns(10)
		poolConnection.SetMaxOpenConns(100)
		poolConnection.SetConnMaxLifetime(time.Hour)
		return db
	}
}
