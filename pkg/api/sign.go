package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"go1f/pkg/conf"
	"log/slog"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type ResponseSign struct {
	Password string `json:"password"`
}

type ResponseToken struct {
	Token string `json:"token"`
}

func signinHandler(w http.ResponseWriter, r *http.Request) {

}

func auth(next http.HandlerFunc, cfg *conf.Configuration) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(cfg.Password) > 0 {
			var token string
			cookie, err := r.Cookie("token")
			if err == nil {
				token = cookie.Value
			}
			jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
				return []byte(cfg.Password), nil
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

func signinWrapper(next http.HandlerFunc, cfg *conf.Configuration) http.HandlerFunc {
	// Обёртка нужна, чтоб передать параметром cfg.password, без обёртки не нашла как передать
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user ResponseSign
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			writeJson(w, ResponseErr{Error: "ошибка передачи данных"}, http.StatusBadRequest)
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &user); err != nil {
			writeJson(w, ResponseErr{Error: "ошибка десериализации JSON"}, http.StatusBadRequest)
			return
		}
		if len(cfg.Password) > 0 {
			if cfg.Password == user.Password {
				h := sha256.New()
				h.Write([]byte(cfg.Password))
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
				writeJson(w, token, http.StatusOK)
				return
			}
			var resErr ResponseErr
			http.Error(w, "неверный пароль", http.StatusUnauthorized)
			resErr.Error = "неверный пароль"
			slog.Error("Sign:", "err", resErr.Error)
			writeJson(w, resErr, http.StatusOK)
		}
		next(w, r)
	})
}
