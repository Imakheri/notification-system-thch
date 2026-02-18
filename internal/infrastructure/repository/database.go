package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/imakheri/notifications-thch/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	DatabaseName       string
	DatabaseUser       string
	DataBasePassword   string
	DatabasePath       string
	DatabasePort       string
	DatabaseConnection *gorm.DB
}

func NewDatabase(cfg *config.Config) *Database {
	db := &Database{
		DatabaseName:     cfg.DatabaseName,
		DatabaseUser:     cfg.DatabaseUser,
		DataBasePassword: cfg.DataBasePassword,
		DatabasePath:     cfg.DatabasePath,
		DatabasePort:     cfg.DatabasePort,
	}
	db.DatabaseConnection = db.Connection()
	return db
}

func (db *Database) Connection() *gorm.DB {

	if db.DatabaseConnection != nil {
		return db.DatabaseConnection
	}

	var dns = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v%v",
		db.DatabaseUser,
		db.DataBasePassword,
		db.DatabasePath,
		db.DatabasePort,
		db.DatabaseName,
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
