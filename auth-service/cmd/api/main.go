package main

import (
	"log"

	"github.com/BernardN38/social-stream-backend/auth-service/internal/application"
)

func main() {
	log.Fatal(application.NewApp().Run())
}
