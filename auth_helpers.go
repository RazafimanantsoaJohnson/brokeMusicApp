package main

import (
	"net/http"

	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/auth"
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
