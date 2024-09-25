package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/BernardN38/social-stream-backend/auth-service/internal/models"
	"github.com/BernardN38/social-stream-backend/auth-service/internal/service"
	"github.com/go-playground/validator/v10"
)

type HandlerInterface interface {
	CheckHealth(w http.ResponseWriter, r *http.Request)
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
}
type Handler struct {
	service *service.Service
}

func NewHandler() *Handler {
	return &Handler{}
}
func (h *Handler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("auth service up and running"))
	// w.WriteHeader(http.StatusOK)
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var registerUserPayload models.RegisterUserPayload
	err := json.NewDecoder(r.Body).Decode(&registerUserPayload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = validator.New().Struct(registerUserPayload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.service.RegisterUser(registerUserPayload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginUserPayload models.LoginUserPayload
	err := json.NewDecoder(r.Body).Decode(&loginUserPayload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = validator.New().Struct(loginUserPayload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId, err := h.service.LoginUser(loginUserPayload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("%v", userId)))
}
