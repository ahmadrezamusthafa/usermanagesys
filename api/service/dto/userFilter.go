package dto

type UserFilterDto struct {
	FullName  *string `json:"full_name"`
	UserEmail *string `json:"user_email"`
	MSISDN    *string `json:"msisdn"`
	Status    *int    `json:"status"`
}
