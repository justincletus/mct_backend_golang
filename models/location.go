package models

type Location struct {
	Id       uint   `db:"id" json:"id"`
	Latitude string `db:"latitude" json:"lat"`
	Longtude string `db:"longtude" json:"lng"`
}
