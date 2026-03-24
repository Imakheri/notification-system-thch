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
		DatabasePath:     cfg.DatabaseHost,
		DatabasePort:     cfg.DatabasePort,
	}
	db.DatabaseConnection = db.Connect()
	return db
}

func (db *Database) Connect() *gorm.DB {

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

	var gormDB *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		log.Printf("Trying to connect to database (Try %d/10)", i)

		gormDB, err = gorm.Open(mysql.Open(dns), &gorm.Config{})

		if err == nil {
			log.Println("Successfully connected to database")
			break
		}

		log.Printf("Error connecting: %v, Retrying in 3 seconds...", err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("could not connect to database after 10 attempts: %v", err)
	}

	poolConnection, err := gormDB.DB()
	if err != nil {
		log.Fatal(err)
	}

	poolConnection.SetMaxIdleConns(10)
	poolConnection.SetMaxOpenConns(100)
	poolConnection.SetConnMaxLifetime(time.Hour)

	db.DatabaseConnection = gormDB
	return gormDB

}
