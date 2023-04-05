package models

import "gorm.io/gorm"

type Location struct {
	gorm.Model

	Id        uint   `db:"id" json:"id"`
	Latitude  string `db:"latitude" json:"latitude"`
	Longitude string `db:"longitude" json:"longitude"`
}
