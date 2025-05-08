package models

type CheckIINRequest struct {
	IIN string `json:"iin"`
}

type CheckIINResponse struct {
	Valid  bool   `json:"valid"`
	Date   string `json:"date,omitempty"`
	Gender string `json:"gender,omitempty"`
	Error  string `json:"error,omitempty"`
}

type CreatePersonRequest struct {
	IIN   string `json:"iin"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type CreatePersonResponse struct {
	Success bool   `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
}

type PersonDTO struct {
	IIN   string `json:"iin"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}
