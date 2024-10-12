package router

import (
	"time"

	"github.com/BernardN38/social-stream-backend/user-service/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	R *chi.Mux
}

func NewRouter(h handler.HandlerInterface) *Router {
	r := chi.NewRouter()

	//middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/api/v1/user/health", h.CheckHealth)
	return &Router{
		R: r,
	}
}
