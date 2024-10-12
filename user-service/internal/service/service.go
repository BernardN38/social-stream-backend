package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/BernardN38/social-stream-backend/user-service/internal/models"
	users_sql "github.com/BernardN38/social-stream-backend/user-service/internal/sqlc/users"
)

type ServiceInterface interface {
	CreateUser(context.Context, models.CreateUserPayload) error
	GetUser(context.Context, models.CreateUserPayload) error
	UpdateUser(context.Context, models.CreateUserPayload) error
	DeleteUser(context.Context, models.CreateUserPayload) error
}
type Service struct {
	db          *sql.DB
	userQueries users_sql.Queries
}

func NewService(db *sql.DB) *Service {
	userQueries := users_sql.New(db)
	return &Service{
		db:          db,
		userQueries: *userQueries,
	}
}

func (s *Service) CreateUser(ctx context.Context, payload models.CreateUserPayload) error {
	timoutCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	successCh := make(chan struct{})
	errCh := make(chan error)
	go func() {
		err := s.userQueries.CreateUser(timoutCtx, users_sql.CreateUserParams{
			Username:  payload.Username,
			Email:     payload.Email,
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Dob:       payload.Dob,
		})
		if err != nil {
			errCh <- err
			return
		}
		successCh <- struct{}{}
	}()
	select {
	case <-successCh:
		return nil
	case err := <-errCh:
		return err
	case <-timoutCtx.Done():
		return timoutCtx.Err()
	}
}

func (s *Service) GetUser(context.Context, models.CreateUserPayload) error {
	return nil
}
func (s *Service) UpdateUser(context.Context, models.CreateUserPayload) error {
	return nil
}
func (s *Service) DeleteUser(context.Context, models.CreateUserPayload) error {
	return nil
}
