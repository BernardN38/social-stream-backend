package handler

import (
	"net/http"
	"time"
)

func CreateJWTCookie(jwtToken string) *http.Cookie {
	return &http.Cookie{
		Name:     "jwt",
		Value:    jwtToken,
		Path:     "/",
		Domain:   "localhost",
		Expires:  time.Now().Add(time.Minute * 30),
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
}
