package middlewares

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
)

var (
	cookieUserName = "UserID"
	userID         []byte
	secretKey      []byte
)

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func checkCookie(r *http.Request) (bool, error) {

	cookie, err := r.Cookie(cookieUserName)

	if errors.Is(err, http.ErrNoCookie) {
		return false, nil
	}

	data, err := hex.DecodeString(cookie.Value)
	if err != nil {
		return false, err
	}

	h := hmac.New(sha256.New, secretKey)
	h.Write(data[:4])

	sign := h.Sum(nil)

	if hmac.Equal(sign, data[4:]) {
		return true, nil
	}

	return false, nil
}

func CookieHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		isCookieCorrect, err := checkCookie(r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !isCookieCorrect {

			if len(userID) == 0 {
				userID, err = generateRandom(4)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			if len(secretKey) == 0 {
				secretKey, err = generateRandom(16)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			h := hmac.New(sha256.New, secretKey)
			h.Write(userID)

			dst := h.Sum(nil)

			http.SetCookie(w, &http.Cookie{
				Name:  cookieUserName,
				Value: fmt.Sprintf("%x", userID) + fmt.Sprintf("%x", dst),
			})
		}

		next.ServeHTTP(w, r)
	})
}
