package models

type RegisterUserPayload struct {
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Username  string `json:"username" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	DOB       string `json:"dob" validate:"required"`
	ImageID   string `json:"imageID"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
