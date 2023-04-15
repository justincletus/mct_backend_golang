package models

import "gorm.io/gorm"

type FeedBack struct {
	gorm.Model

	Id          uint   `db:"id" json:"id" gorm:"primaryKey; autoIncrement:true"`
	Title       string `db:"title" json:"title" binding:"required title"`
	Description string `db:"description" json:"description" gorm:"type:text"`
	Uid         uint   `db:"uid" json:"uid" gorm:"index"`
	Username    string `db:"username" json:"username"`
	Address     string `db:"address" json:"address" gorm:"default:null"`
	User        User   `json:"user" gorm:"foreignKey:Uid"`
}
