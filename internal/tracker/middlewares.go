package tracker

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type Middleware struct {
	l *slog.Logger
}

func NewMiddleware(l *slog.Logger) *Middleware {
	return &Middleware{
		l: l,
	}
}

type LoggerCtxKey struct{}

func (m *Middleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := m.l.With("request_id", uuid.New())

		l.Info("incoming request", "method", r.Method, "url", r.URL.String(), "from", r.RemoteAddr)

		ctx := context.WithValue(r.Context(), LoggerCtxKey{}, l)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
