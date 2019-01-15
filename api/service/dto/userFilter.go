package dto

type UserFilterDto struct {
	FullName *string `json:"full_name"`
	Status   *int    `json:"status"`
}
