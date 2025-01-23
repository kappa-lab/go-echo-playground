package logger

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type loggerKey struct{}

func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(loggerKey{}).(*zap.Logger)
	if ok {
		return logger
	}
	return zap.NewExample().With(
		zap.Bool("is_example", true),
		zap.String("warn", "contextにLoggerがセットされていません"),
	)
}

func LoggerMiddleware(l *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			ctx = WithContext(ctx, l)
			c.SetRequest(c.Request().WithContext(ctx))

			err := next(c)

			tl := FromContext(c.Request().Context()).
				With(zap.String("method", c.Request().Method)).
				With(zap.String("path", c.Request().RequestURI)).
				With(zap.Int("status", c.Response().Status))

			if err != nil {
				c.Error(err)
				tl.WithOptions(zap.AddCallerSkip(1)).Error("NG")
			}

			tl.Info("OK")

			return nil
		}
	}
}
