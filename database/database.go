package database

import (
	"fmt"

	"github.com/justincletus/cms/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connetion() error {
	data, err := config.Config()
	if err != nil {
		return err
	}

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", data["username"], data["password"], data["host"], data["port"], data["database"])

	dataSource += "?charset=utf8mb4&parseTime=True&loc=Local"
	connection, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = connection
	// DB.AutoMigrate(&models.Report{})
	// DB.AutoMigrate(&models.User{})
	// DB.AutoMigrate(&models.Manager{})
	// DB.AutoMigrate(&models.TeamMem{})

	return nil

}
