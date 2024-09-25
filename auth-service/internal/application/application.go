package application

import (
	"net/http"

	"github.com/BernardN38/social-stream-backend/auth-service/internal/handler"
	"github.com/BernardN38/social-stream-backend/auth-service/internal/router"
)

type App struct {
	router *router.Router
}

func NewApp() *App {
	hanlder := handler.NewHandler()
	router := router.NewRouter(hanlder)
	return &App{
		router: router,
	}
}

func (a *App) Run() error {
	return http.ListenAndServe(":8080", a.router.R)
}
