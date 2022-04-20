package handler

import "net/http"

// CustomMiddleware custom user middleware
func CustomMiddleware(_ interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			/*
				use _ param to do something
			*/
			next.ServeHTTP(writer, request)
		})
	}
}
