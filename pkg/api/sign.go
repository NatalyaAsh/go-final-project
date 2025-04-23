package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

type ResponseSign struct {
	Password string `json:"password"`
}

type ResponseToken struct {
	Token string `json:"token"`
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var user ResponseSign
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		writeJson(w, ResponseErr{Error: "ошибка передачи данных"})
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		writeJson(w, ResponseErr{Error: "ошибка десериализации JSON"})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if len(pass) > 0 {
		if pass == user.Password {
			h := sha256.New()
			h.Write([]byte(pass))
			hashSum := h.Sum(nil)

			claims := jwt.MapClaims{
				"hashSum": hashSum,
			}
			jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			var token ResponseToken
			token.Token, err = jwtToken.SignedString([]byte(user.Password))
			if err != nil {
				slog.Error("failed to sign jwt:", "err", err.Error())
			}
			writeJson(w, token)
			return
		}
		var resErr ResponseErr
		http.Error(w, "неверный пароль", http.StatusUnauthorized)
		resErr.Error = "неверный пароль"
		slog.Error("Sign:", "err", resErr.Error)
		writeJson(w, resErr)
	}
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var token string
			cookie, err := r.Cookie("token")
			if err == nil {
				token = cookie.Value
			}
			jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				return []byte(pass), nil
			})
			if err != nil {
				slog.Error("failed to parse token", "err", err)
				return
			}
			if !jwtToken.Valid {
				slog.Error("Токен невалиден")
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
