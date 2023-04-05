package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Email    string    `db:"email" json:"email" binding:"required email" gorm:"unique; not null"`
	Fullname string    `db:"fullname" json:"fullname"`
	Mobile   string    `db:"mobile" json:"mobile"`
	Password []byte    `db:"password" json:"-"`
	Location *Location `json:",omitempty" gorm:"foreignKey:id"`
}
