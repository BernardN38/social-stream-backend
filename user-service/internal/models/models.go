package models

type CreateUserPayload struct {
	Username  string `json:"username" vaildate:"required"`
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Dob       string `json:"dob" validate:"required"`
}
