package handler

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var compressible = map[string]string{
	"text/html":                "",
	"text/css":                 "",
	"text/plain":               "",
	"text/javascript":          "",
	"application/javascript":   "",
	"application/x-javascript": "",
	"application/json":         "",
	"application/atom+xml":     "",
}

type gzipResponseWriter struct {
	gz io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {

	if _, ok := compressible[w.Header().Get("Content-Type")]; ok {
		w.Header().Del("Content-Length")
		return w.gz.Write(b)
	}

	return w.ResponseWriter.Write(b)
}

func CompressMiddleware(_ interface{}) func(http.Handler) http.Handler {
	gzWriter := gzip.NewWriter(nil)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

			// gzipped request, try to process
			var req http.Request
			if request.Header.Get(`Content-Encoding`) == "gzip" {
				req = *request
				gzReader, err := gzip.NewReader(req.Body)
				if err != nil {
					http.Error(writer, err.Error(), http.StatusInternalServerError)
					return
				}
				request.Body = gzReader
				defer gzReader.Close()
			}

			// now check if client is able to receive compressed content
			if !strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(writer, request)
				return
			}

			writer.Header().Set("Content-Encoding", "gzip")
			gzWriter.Reset(writer)
			defer gzWriter.Close()

			gzw := gzipResponseWriter{
				gz:             gzWriter,
				ResponseWriter: writer,
			}

			next.ServeHTTP(gzw, request)

		})
	}
}
