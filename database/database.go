package database

import (
	"github.com/diyor200/go-fintech/helpers"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var database, err = gorm.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=2001 dbname=bankapp sslmode=disable")
	helpers.HandleErr(err)
	database.DB().SetMaxIdleConns(20)
	database.DB().SetMaxOpenConns(200)
	DB = database
}
