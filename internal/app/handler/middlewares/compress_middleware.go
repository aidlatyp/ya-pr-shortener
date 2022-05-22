package middlewares

import (
	"compress/gzip"
	"io"
	"log"
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
	http.ResponseWriter
	gzipWriter io.Writer
}

// canCompress checks if the Content-Type header represent a compressible type of content.
func (grw *gzipResponseWriter) canCompress() bool {
	_, ok := compressible[grw.Header().Get("Content-Type")]
	return ok
}

func (grw *gzipResponseWriter) Write(b []byte) (int, error) {
	if grw.canCompress() {
		if closer, ok := (grw.gzipWriter).(*gzip.Writer); ok {
			defer func() {
				err := closer.Close()
				if err != nil {
					log.Printf("error while closing gzip writer, %v", err)
				}
			}()
		}
		// compressed
		return grw.gzipWriter.Write(b)
	}

	// ordinary
	return grw.ResponseWriter.Write(b)
}

func (grw *gzipResponseWriter) WriteHeader(statusCode int) {
	if grw.canCompress() {
		grw.Header().Set("Content-Encoding", "gzip")
	}
	grw.ResponseWriter.WriteHeader(statusCode)
}

func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		// decompress Request
		var req http.Request
		if request.Header.Get(`Content-Encoding`) == "gzip" {
			req = *request
			gzReader, err := gzip.NewReader(req.Body)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
				return
			}
			request.Body = gzReader

			defer func() {
				err := gzReader.Close()
				if err != nil {
					log.Printf("error while decompress request boby, %v", err)
					writer.WriteHeader(400)
				}
			}()

		}
		if !strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(writer, request)
			return
		}

		// compress Response
		// naive realization
		gzWriter := gzip.NewWriter(writer)
		gzw := gzipResponseWriter{
			gzipWriter:     gzWriter,
			ResponseWriter: writer,
		}

		next.ServeHTTP(&gzw, request)
	})
}
