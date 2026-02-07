package repository

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Database() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal()
	}
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
