package handler

import (
	"time"

	"github.com/go-chi/jwtauth/v5"
)

type TokenManager struct {
	jwtAuth *jwtauth.JWTAuth
}

func NewTokenManger(jwtAuth *jwtauth.JWTAuth) *TokenManager {
	return &TokenManager{
		jwtAuth: jwtAuth,
	}
}
func (tm *TokenManager) CreateToken(userId int) (string, error) {
	_, tokenString, err := tm.jwtAuth.Encode(map[string]interface{}{"user_id": userId, "exp": time.Now().Add(time.Minute * 30)})
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
