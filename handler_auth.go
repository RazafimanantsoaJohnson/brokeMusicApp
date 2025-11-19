package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/auth"
	"github.com/RazafimanantsoaJohnson/brokeMusicApp/internal/database"
)

type SignupBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SigninBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	AuthType string `json:"auth_type"`
	Token    string `json:"token"`
}

type SigninResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *ApiConfig) HandleSignup(w http.ResponseWriter, r *http.Request) {
	signUpData := SignupBody{}
	defer r.Body.Close()
	signUpBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("unable to read request body"))
		return
	}
	err = json.Unmarshal(signUpBody, &signUpData)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	if signUpData.Email == "" || signUpData.Password == "" {
		w.WriteHeader(400)
		w.Write([]byte("non-valid request body"))
		return
	}

	hashedPassword, err := auth.HashPassword(signUpData.Password)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("unable to encrypt password"))
		return
	}

	newUser, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		Email:    signUpData.Email,
		Password: hashedPassword,
		Authtype: AuthTypes[0],
	})

	if err != nil {
		if strings.Contains(err.Error(), duplicateUserError) {
			w.WriteHeader(400)
			w.Write([]byte("the email is already used by another user"))
			return
		}
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	response := UserResponse{
		Id:       newUser.ID.String(),
		Username: newUser.Username.String,
		Email:    newUser.Email,
		AuthType: newUser.Authtype,
	}
	returnJson(w, response)
}

func (cfg *ApiConfig) HandleSignin(w http.ResponseWriter, r *http.Request) {
	signinData := SigninBody{}
	defer r.Body.Close()
	signinBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("signin body non-valid"))
		return
	}
	err = json.Unmarshal(signinBody, &signinData)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	dbUser, err := cfg.db.FetchUserByEmail(context.Background(), signinData.Email)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	err = auth.CheckPasswordHash(signinData.Password, dbUser.Password)
	if err != nil {
		w.WriteHeader(403)
		w.Write([]byte(err.Error()))
		return
	}
	token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, 1*time.Hour)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	refreshToken, err := createRefreshToken(cfg, dbUser.ID)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	response := SigninResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}

	returnJson(w, response)
}

func (cfg *ApiConfig) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("unauthorized user"))
		return
	}

	dbToken, err := cfg.db.GetTokenById(context.Background(), refreshToken)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte("unauthorized user"))
		return
	}

	if dbToken.Revokedat.Valid {
		w.WriteHeader(401)
		w.Write([]byte("unauthorized user"))
		return
	}

	newAccessToken, err := auth.MakeJWT(dbToken.Userid, cfg.jwtSecret, 1*time.Hour)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("unable to generate access token"))
		return
	}

	response := SigninResponse{
		Token: newAccessToken,
	}

	returnJson(w, response)
}
