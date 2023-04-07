package models

import "gorm.io/gorm"

type Location struct {
	gorm.Model

	Id        uint   `db:"id" json:"id"`
	Latitude  string `db:"latitude" json:"latitude"`
	Longitude string `db:"longitude" json:"longitude"`
	Address   string `db:"address" json:"address"`
	City      string `db:"city" json:"city"`
	State     string `db:"state" json:"state"`
	PinCode   string `db:"pin_code" json:"pin_code"`
	Country   string `db:"country" json:"country"`
}
