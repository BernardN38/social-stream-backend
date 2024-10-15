package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/BernardN38/social-stream-backend/post-service/internal/messaging"
	"github.com/BernardN38/social-stream-backend/post-service/internal/models"
	posts_sql "github.com/BernardN38/social-stream-backend/post-service/internal/sqlc/posts"
)

type ServiceInterface interface {
	CreatePost(ctx context.Context, payload models.CreatePostPayload) error
}
type Service struct {
	db              *sql.DB
	rabbitmeEmitter messaging.MessageEmitter
	postQueries     posts_sql.Queries
}

func NewService(db *sql.DB, rabbitmqEmitter *messaging.RabbitmqEmitter) *Service {
	userQueries := posts_sql.New(db)
	return &Service{
		db:              db,
		rabbitmeEmitter: rabbitmqEmitter,
		postQueries:     *userQueries,
	}
}

func (s *Service) CreatePost(ctx context.Context, payload models.CreatePostPayload) error {
	timoutCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	successCh := make(chan struct{})
	errCh := make(chan error)
	go func() {
		err := s.postQueries.CreatePost(timoutCtx, posts_sql.CreatePostParams{
			OwnerID: int32(payload.OwnerID),
			Body:    payload.Body,
			MediaID: sql.NullString{
				String: payload.MediaID,
				Valid:  len(payload.MediaID) > 0,
			},
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

// func (s *Service) GetUser(context.Context, models.CreateUserPayload) error {
// 	return nil
// }
// func (s *Service) UpdateUser(context.Context, models.CreateUserPayload) error {
// 	return nil
// }
// func (s *Service) DeleteUser(context.Context, models.CreateUserPayload) error {
// 	return nil
// }
