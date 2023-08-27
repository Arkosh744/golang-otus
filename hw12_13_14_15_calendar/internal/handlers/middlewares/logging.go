package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/Arkosh744/otus-go/hw12_13_14_15_calendar/internal/log"
)

type StatusResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (s *StatusResponseWriter) WriteHeader(status int) {
	s.Status = status
	s.ResponseWriter.WriteHeader(status)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		sw := &StatusResponseWriter{ResponseWriter: w}

		next.ServeHTTP(sw, r)

		elapsed := time.Since(timeStart)

		ip := strings.SplitN(r.RemoteAddr, ":", 2)[0]

		log.Infof(
			"%s %s %s HTTP/%d.%d %d %dms \"%s\"",
			ip,
			r.Method,
			r.RequestURI,
			r.ProtoMajor,
			r.ProtoMinor,
			sw.Status,
			elapsed.Milliseconds(),
			r.UserAgent(),
		)
	})
}
