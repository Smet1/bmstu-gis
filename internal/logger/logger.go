package logger

import (
	"context"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type ctxlog struct{}

// WithLogger put logger to context
func WithLogger(ctx context.Context, logger *logrus.Logger) context.Context {
	return context.WithValue(ctx, ctxlog{}, logger)
}

// GetLogger get logger from context
func GetLogger(ctx context.Context) *logrus.Logger {
	l, ok := ctx.Value(ctxlog{}).(logrus.Logger)
	if !ok {
		l = *logrus.New()
		l.SetOutput(os.Stdout)
		l.SetLevel(logrus.InfoLevel)
	}
	return &l
}

// GetLoggerMiddleware get middleware for router with logger
func GetLoggerMiddleware(log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			ctx := WithLogger(req.Context(), log)

			next.ServeHTTP(res, req.WithContext(ctx))
		})
	}
}
