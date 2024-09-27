package messaging

type CreateUserMessage struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Dob       string `json:"dob"`
}
