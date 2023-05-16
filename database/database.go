package database

import (
	"github.com/justincletus/cms/config"
	"github.com/justincletus/cms/models"
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
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.FeedBack{})

	return nil

}
