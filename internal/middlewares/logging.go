package middlewares

import (
	"net/http"
	"time"

	"github.com/Soliard/gophermart/internal/logger"
	"github.com/google/uuid"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData: &responseData{
				status: 0,
				size:   0,
			},
		}
		requestID := uuid.NewString()
		log := logger.FromContext(ctx).With(
			logger.F.String("request id", requestID),
		)

		log.Info("request info",
			logger.F.String("url", r.URL.String()),
			logger.F.String("method", r.Method),
		)

		ctx = logger.WithContext(ctx, log)
		next.ServeHTTP(&lw, r.WithContext(ctx))
		duration := time.Since(start)

		log.Info("response info",
			logger.F.Duration("duration", duration),
			logger.F.Int("size", lw.responseData.size),
			logger.F.Int("status", lw.responseData.status),
		)
	})
}
