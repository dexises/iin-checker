package models

// CheckIINRequest represents the JSON payload for IIN validation.
type CheckIINRequest struct {
	IIN string `json:"iin"`
}

// CheckIINResponse represents the JSON response for IIN validation.
type CheckIINResponse struct {
	Valid  bool   `json:"valid"`
	Date   string `json:"date,omitempty"`
	Gender string `json:"gender,omitempty"`
	Error  string `json:"error,omitempty"`
}

// CreatePersonRequest represents payload to create a person record.
type CreatePersonRequest struct {
	IIN   string `json:"iin"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// CreatePersonResponse represents JSON response after creating a person.
type CreatePersonResponse struct {
	Success bool   `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`
}

// PersonDTO represents a person record.
type PersonDTO struct {
	IIN   string `json:"iin"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}
