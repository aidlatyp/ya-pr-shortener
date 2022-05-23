package middlewares

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/aidlatyp/ya-pr-shortener/internal/util"
)

const secret = "secret"

type key int

const (
	UserIDCtxKey key = iota
)

func registerUser(writer http.ResponseWriter) []byte {
	userID := util.GenerateUserID()
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(userID)
	signed := h.Sum(nil)

	c := http.Cookie{
		Name:   "user_id",
		Value:  hex.EncodeToString(append(userID, signed...)),
		MaxAge: 3600 * 24,
	}
	http.SetCookie(writer, &c)
	return userID
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		var userID []byte
		signedCookie, err := request.Cookie("user_id")

		if err != nil {
			// Register user silently
			userID = registerUser(writer)
		} else {

			cookieBytes, err := hex.DecodeString(signedCookie.Value)
			if err != nil {
				log.Println(err)
			}
			userID = cookieBytes[:6]
			incomeSign := cookieBytes[6:]
			h := hmac.New(sha256.New, []byte(secret))
			h.Write(userID)
			controlSign := h.Sum(nil)

			if !hmac.Equal(incomeSign, controlSign) {
				// Register user silently
				userID = registerUser(writer)
			}
		}

		userCtx := context.WithValue(request.Context(), UserIDCtxKey, string(userID))
		request = request.WithContext(userCtx)
		next.ServeHTTP(writer, request)
	})
}
