package dto

type APIResultDto struct {
	Result  interface{} `json:"result"`
	Success bool        `json:"success"`
	Error   *string     `json:"error""`
}
