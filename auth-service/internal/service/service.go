package service

import (
	"github.com/BernardN38/social-stream-backend/auth-service/internal/models"
)

type Service struct {
}

func (s *Service) RegisterUser(payload models.RegisterUserPayload) error {
	return nil
}

func (s *Service) LoginUser(payload models.LoginUserPayload) (int, error) {
	return 1, nil
}
