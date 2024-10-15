package router

import (
	"time"

	"github.com/BernardN38/social-stream-backend/post-service/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

type Router struct {
	R *chi.Mux
}

func NewRouter(h handler.HandlerInterface, tokenAuth *jwtauth.JWTAuth) *Router {
	r := chi.NewRouter()

	//middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(jwtauth.Authenticator(tokenAuth))
	r.Get("/api/v1/post/health", h.CheckHealth)
	r.Post("/api/v1/post", h.CreatePost)
	return &Router{
		R: r,
	}
}
