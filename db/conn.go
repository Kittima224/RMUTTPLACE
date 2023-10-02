package db

import (
	"RmuttPlace/model"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Conn *gorm.DB
var err error

func InitDB() {
	dsn := os.Getenv("MYSQL_DNS")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic("failed to connect database")
	}

	Conn = db
}

func Migrate() {
	Conn.AutoMigrate(
		&model.User{},
		&model.Store{},
		&model.Product{},
		&model.Review{},
		&model.Shipment{},
		&model.Cart{},
		&model.Favorite{},
		&model.Admin{},
		&model.Order{},
		&model.OrderItem{},
		&model.Category{},
	)
}
