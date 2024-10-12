package service

type UserNotFoundError struct {
	message string
}
type UnauthorizedError struct {
	message string
}

func (e UserNotFoundError) Error() string {
	return e.message
}

func (e UnauthorizedError) Error() string {
	return e.message
}
