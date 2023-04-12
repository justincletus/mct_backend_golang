package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Id       uint   `db:"id" json:"id" gorm:"primaryKey;autoIncrement:true"`
	Email    string `db:"email" json:"email" binding:"required email" gorm:"unique; not null"`
	Username string `db:"username" json:"username" gorm:"unique"`
	Fullname string `db:"fullname" json:"fullname"`
	Mobile   string `db:"mobile" json:"mobile"`
	Password []byte `db:"password" json:"-"`
	//Location *Location `json:",omitempty" gorm:"foreignKey:Location"`
}
