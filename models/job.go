package models

import "time"

type Job struct {
	Id        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UserId    int       `db:"user_id" json:"user_id"`
	User      User
}
