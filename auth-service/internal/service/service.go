package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/BernardN38/ebuy-server/authentication-service/messaging"
	"github.com/BernardN38/social-stream-backend/auth-service/internal/models"

	users_sql "github.com/BernardN38/social-stream-backend/auth-service/sqlc/users"
	"golang.org/x/crypto/bcrypt"
)

type ServiceInterface interface {
	RegisterUser(ctx context.Context, payload models.RegisterUserPayload) error
	LoginUser(ctx context.Context, payload models.LoginUserPayload) (int, error)
}
type Service struct {
	db              *sql.DB
	userQueries     users_sql.Queries
	rabbitmqEmitter messaging.MessageEmitter
}

func NewService(db *sql.DB, rabbitmqEmitter messaging.MessageEmitter) *Service {
	userQueries := users_sql.New(db)
	return &Service{
		db:              db,
		userQueries:     *userQueries,
		rabbitmqEmitter: rabbitmqEmitter,
	}
}
func (s *Service) RegisterUser(ctx context.Context, payload models.RegisterUserPayload) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5000*time.Millisecond)
	defer cancel()
	successCh := make(chan struct{})
	errorCh := make(chan error)
	go func() {
		encodedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 12)
		if err != nil {
			errorCh <- err
			return
		}
		err = s.userQueries.CreateUser(timeoutCtx, users_sql.CreateUserParams{
			Username:        payload.Username,
			Email:           payload.Email,
			EncodedPassword: string(encodedPassword),
		})
		if err != nil {
			errorCh <- err
			return
		}
		msg, err := json.Marshal(messaging.CreateUserMessage{
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Username:  payload.Username,
			Email:     payload.Email,
			Dob:       payload.DOB,
		})
		if err != nil {
			errorCh <- err
			return
		}
		err = s.rabbitmqEmitter.SendMessage(timeoutCtx, msg, "user_events", "user.created", "user.created")
		if err != nil {
			errorCh <- err
			return
		}
		successCh <- struct{}{}
	}()
	select {
	case <-successCh:
		return nil
	case err := <-errorCh:
		return err
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}

func (s *Service) LoginUser(ctx context.Context, payload models.LoginUserPayload) (int, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5000*time.Millisecond)
	defer cancel()
	successCh := make(chan int32)
	errorCh := make(chan error)
	go func() {
		user, err := s.userQueries.GetUserPassword(timeoutCtx, payload.Email)
		if err != nil {
			errorCh <- UserNotFoundError{
				message: "user not found",
			}
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.EncodedPassword), []byte(payload.Password))
		if err != nil {
			errorCh <- UnauthorizedError{
				message: "could not authenticate user",
			}
			return
		}
		successCh <- user.ID
	}()
	select {
	case userId := <-successCh:
		return int(userId), nil
	case err := <-errorCh:
		return 0, err
	case <-timeoutCtx.Done():
		return 0, timeoutCtx.Err()
	}
}
