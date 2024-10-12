package router

import (
	"time"

	"github.com/BernardN38/social-stream-backend/auth-service/internal/handler"
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
	r.Use(middleware.RequestSize(50 * 1024 * 1024))
	r.Get("/api/v1/auth/health", h.CheckHealth)
	r.Post("/api/v1/auth/register", h.RegisterUser)
	r.Post("/api/v1/auth/login", h.LoginUser)
	return &Router{
		R: r,
	}
}
