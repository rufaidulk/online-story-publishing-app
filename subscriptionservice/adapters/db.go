package adapters

import (
	"log"
	"subscriptionservice/helper"
	"sync"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var dbConnector sync.Once

func GetDbHandle(ctx echo.Context) *gorm.DB {
	dbConnector.Do(func() {
		mode := ctx.Get("mode")
		if mode == "testing" {
			initTestDb()
		} else {
			initDb()
		}
	})

	return db
}

func initDb() {
	var err error
	dsn := helper.GetDsn()
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Error on db connection")
		log.Fatal(err)
	} else {
		log.Println("Connected to database")
	}
}

func initTestDb() {
	var err error
	dsn := helper.GetTestDbDsn()
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Error on test db connection")
		log.Fatal(err)
	} else {
		log.Println("Connected to test database")
	}
}
