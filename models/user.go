package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Id            uint    `db:"id" json:"id" gorm:"primaryKey;autoIncrement:true"`
	Email         string  `db:"email" json:"email" binding:"required email" gorm:"unique; not null"`
	Username      string  `db:"username" json:"username" gorm:"unique"`
	Fullname      string  `db:"fullname" json:"fullname"`
	Mobile        string  `db:"mobile" json:"mobile"`
	Password      []byte  `db:"password" json:"-"`
	Role          string  `db:"role" json:"role" gorm:"default:null"`
	Code          string  `db:"code" json:"code" gorm:"default:null"`
	EmailVerified bool    `db:"email_verified" json:"email_verified" gorm:"default:false"`
	Manager       Manager `gorm:"default:null; constraint:OnDelete:CASCADE"`
	TeamMems      []TeamMem
}

type Manager struct {
	gorm.Model

	Id     uint `db:"id" json:"id" gorm:"primaryKey; autoIncrement:true"`
	UserId uint
}

type TeamMem struct {
	gorm.Model
	Id      uint   `db:"id" json:"id" gorm:"primaryKey; autoIncrement:true"`
	Title   string `db:"title" json:"title"`
	UserId  uint   `json:"user_id" gorm:"foreignKey:UserId"`
	User    User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Members string `db:"members" json:"members"`
}
