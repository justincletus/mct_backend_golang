package models

type User struct {
	Id            int    `db:"id" json:"id"`
	Email         string `db:"email" json:"email"`
	Username      string `db:"username" json:"username"`
	Fullname      string `db:"fullname" json:"fullname"`
	Mobile        string `db:"mobile" json:"mobile"`
	Password      []byte `db:"password" json:"-"`
	Role          string `db:"role" json:"role"`
	Code          string `db:"code" json:"code"`
	EmailVerified bool   `db:"email_verified" json:"email_verified"`
}

type Manager struct {
	Id     int `db:"id" json:"id" gorm:"primaryKey; autoIncrement:true"`
	UserId int
}

type TeamMem struct {
	Id              int    `db:"id" json:"id"`
	Title           string `db:"title" json:"title"`
	SubContractor   string `db:"sub_contractor" json:"sub_contractor"`
	ContractorEmail string `db:"contractor_email;omitempty" json:"contractor_email"`
	ClientEmail     string `db:"client_email;omitempty" json:"client_email"`
	Members         string `db:"members;omitempty" json:"members"`
	UserId          int    `json:"user_id"`
	User            User
}

type Member struct {
	Id     int    `db:"id" json:"id"`
	Email  string `db:"email" json:"email"`
	TeamId int    `db:"team_id" json:"team_id"`
}
