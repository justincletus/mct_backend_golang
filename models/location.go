package models

import "gorm.io/gorm"

type Location struct {
	gorm.Model

	Id        uint   `db:"id" json:"id"`
	Latitude  string `db:"latitude" json:"lat"`
	Longitude string `db:"longtude" json:"lng"`
}
