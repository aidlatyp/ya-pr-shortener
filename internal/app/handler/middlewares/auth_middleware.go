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

		// register
		if err != nil {
			userID = registerUser(writer)
		} else {

			cookieBytes, err := hex.DecodeString(signedCookie.Value)
			if err != nil {
				log.Println(err)
			}

			userId := cookieBytes[:6]
			incomeSign := cookieBytes[6:]

			h := hmac.New(sha256.New, []byte(secret))
			h.Write(userId)
			controlSign := h.Sum(nil)

			if hmac.Equal(incomeSign, controlSign) {
				userID = userId
			} else {
				userID = registerUser(writer)
			}
		}
		userCtx := context.WithValue(request.Context(), "userID", string(userID))
		request = request.WithContext(userCtx)
		next.ServeHTTP(writer, request)
	})
}
