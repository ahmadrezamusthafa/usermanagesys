package dbo

import "time"

type TokopediaUser struct {
	UserId       int       `json:"user_id"`
	UserName     string    `json:"user_name"`
	UserEmail    string    `json:"user_email"`
	UserPassword string    `json:"user_pwd"`
	Status       int       `json:"status"`
	FullName     string    `json:"full_name"`
	Sex          int       `json:"sex"`
	BirthDate    time.Time `json:"birth_date"`
	Age          int       `json:"age"`
	Location     string    `json:"location"`
	MSISDN       string    `json:"msisdn"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
}
