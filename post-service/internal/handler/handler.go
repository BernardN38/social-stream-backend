package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BernardN38/social-stream-backend/post-service/internal/models"
	"github.com/BernardN38/social-stream-backend/post-service/internal/service"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
)

type HandlerInterface interface {
	CheckHealth(w http.ResponseWriter, r *http.Request)
	CreatePost(w http.ResponseWriter, r *http.Request)
}
type Handler struct {
	service service.ServiceInterface
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		service: s}
}

func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("post service up and running"))
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	log.Println(claims)
	var payload models.CreatePostPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "error decoding json payload", http.StatusBadRequest)
		log.Println(err)
		return
	}
	err = validator.New().Struct(payload)
	if err != nil {
		http.Error(w, "invalid json payload", http.StatusBadRequest)
		log.Println(err)
		return
	}
	err = h.service.CreatePost(r.Context(), payload)
	if err != nil {
		http.Error(w, "error creating post", http.StatusBadRequest)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
