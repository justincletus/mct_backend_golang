package models

import "time"

type Job struct {
	Id        int       `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	JobId     string    `db:"job_id;omitempty" json:"job_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    int       `db:"user_id" json:"user_id"`
	User      User
}
