package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BernardN38/social-stream-backend/auth-service/internal/models"
	"github.com/BernardN38/social-stream-backend/auth-service/internal/service"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-playground/validator/v10"
)

type HandlerInterface interface {
	CheckHealth(w http.ResponseWriter, r *http.Request)
	RegisterUser(w http.ResponseWriter, r *http.Request)
	LoginUser(w http.ResponseWriter, r *http.Request)
}
type Handler struct {
	Service      service.ServiceInterface
	TokenManager *TokenManager
}

func NewHandler(s *service.Service, jwtAuth *jwtauth.JWTAuth) *Handler {
	tm := NewTokenManger(jwtAuth)
	return &Handler{
		Service:      s,
		TokenManager: tm}
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
	err = h.Service.RegisterUser(r.Context(), registerUserPayload)
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
	userId, err := h.Service.LoginUser(r.Context(), loginUserPayload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokenString, err := h.TokenManager.CreateToken(userId)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	cookie := CreateJWTCookie(tokenString)
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
	// w.Write([]byte(fmt.Sprintf("%v", userId)))
}
