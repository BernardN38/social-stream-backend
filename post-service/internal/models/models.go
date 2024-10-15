package models

type CreatePostPayload struct {
	OwnerID int
	Body    string `json:"body" validate:"required"`
	MediaID string `json:"mediaID"`
}
