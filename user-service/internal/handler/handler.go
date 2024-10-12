package handler

import (
	"net/http"

	"github.com/BernardN38/social-stream-backend/user-service/internal/service"
)

type HandlerInterface interface {
	CheckHealth(w http.ResponseWriter, r *http.Request)
}
type Handler struct {
	Service service.ServiceInterface
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		Service: s}
}

func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("user service up and running"))
}
