package database

import (
	"github.com/justincletus/map-backend/config"
	"github.com/justincletus/map-backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connetion() error {
	dataSource, err := config.Config()
	if err != nil {
		return err
	}

	dataSource += "?charset=utf8mb4&parseTime=True&loc=Local"
	connection, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = connection
	DB.AutoMigrate(&models.Location{})

	return nil

}
