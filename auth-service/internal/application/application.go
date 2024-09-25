package application

import (
	"net/http"

	"github.com/BernardN38/social-stream-backend/auth-service/internal/handler"
	"github.com/BernardN38/social-stream-backend/auth-service/internal/router"
)

type App struct {
	Router *router.Router
}

func NewApp() *App {
	hanlder := handler.NewHandler()
	router := router.NewRouter(hanlder)
	return &App{
		Router: router,
	}
}

func (a *App) Run() error {
	return http.ListenAndServe(":8080", a.Router.R)
}
