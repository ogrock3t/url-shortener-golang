package logger

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

type contextKey string

const RequestID contextKey = "request_id"

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func MiddlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := generateRequestID()

		rw := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		log.Printf(
			"request id = %s %s %s %d %s %s",
			requestID,
			r.Method,
			r.URL.Path,
			rw.status,
			start,
			duration,
		)
	})
}

func generateRequestID() string {
	alphabet := "abcdefghijklmnopqrstuvwxyz" + "0123456789"

	b := make([]byte, 8)
	for i := range b {
		b[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(b)
}
