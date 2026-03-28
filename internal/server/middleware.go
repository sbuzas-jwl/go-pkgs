package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/sbuzas-jwl/go-pkgs/internal/logging"
)

func RequestLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := logger
		if logger == nil {
			log = logging.DefaultLogger()
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			httpLog := log.WithGroup("http").
				With(
					slog.String("method", r.Method),
					slog.String("host", r.Host),
					slog.String("url", r.URL.String()),
					slog.Int("proto_major", r.ProtoMajor),
					slog.Int("proto_minor", r.ProtoMinor),
					slog.String("route", r.Pattern),
					slog.String("user_agent", r.UserAgent()),
					slog.Int64("content_length", r.ContentLength),
				)
			if v := r.Referer(); v != "" {
				httpLog = httpLog.With(slog.String("referer", v))
			}
			lrw := &loggingResponseWriter{ResponseWriter: w}
			next.ServeHTTP(lrw, r)

			httpLog.With(
				slog.Int64("duration", time.Since(start).Microseconds()),
				slog.Int("status_code", lrw.statusCode),
				slog.Int("response_bytes", lrw.bytesCount),
			).Info("")
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	bytesCount int
}

// Write counts the bytes written and calls the underlying ResponseWriter's Write method.
func (w *loggingResponseWriter) Write(p []byte) (int, error) {
	n, err := w.ResponseWriter.Write(p)
	if err == nil {
		w.bytesCount += n // Increment the count with the number of bytes successfully written
	}
	return n, err
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
