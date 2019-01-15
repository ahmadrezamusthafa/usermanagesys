package dbo

import "time"

type TokopediaUser struct {
	UserId       int        `json:"user_id",dbname:"user_id"`
	UserName     *string    `json:"user_name",dbname:"user_name"`
	UserEmail    string     `json:"user_email",dbname:"user_email"`
	UserPassword *string    `json:"user_password",dbname:"user_pwd"`
	Status       int        `json:"status",dbname:"status"`
	FullName     string     `json:"full_name",dbname:"full_name"`
	Sex          *int       `json:"sex",dbname:"sex"`
	BirthDate    *time.Time `json:"birth_date",dbname:"birth_date"`
	Age          *int       `json:"age",dbname:"age"`
	Location     *string    `json:"location",dbname:"location"`
	MSISDN       *string    `json:"msisdn",dbname:"msisdn"`
	CreateTime   *time.Time `json:"create_time",dbname:"create_time"`
	UpdateTime   *time.Time `json:"update_time",dbname:"update_time"`
}
