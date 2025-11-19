package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/auth"
	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) middlewareCheckAuth(next func(*ApiConfig, uuid.UUID, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		receivedToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(UnauthorizedErrorMessage))
			return
		}
		curUserId, err := auth.ValidateJWT(receivedToken, cfg.jwtSecret)
		if err != nil {
			w.WriteHeader(401)
			w.Write([]byte(UnauthorizedErrorMessage))
			return
		}

		next(cfg, curUserId, w, r)
	}
}

func createRefreshToken(cfg *ApiConfig, userId uuid.UUID) (string, error) {
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		return "", err
	}

	_, err = cfg.db.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		Userid:    userId,
		Revokedat: sql.NullTime{Valid: true, Time: time.Now().Add(RefreshTokenDurationInMonth * (30 * 24) * time.Hour)},
	})
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}
