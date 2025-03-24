package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func LogrusMiddleware(logger *logrus.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogMethod:    true,
		LogRemoteIP:  true,
		LogUserAgent: true,
		LogLatency:   true,
		LogError:     true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.WithFields(logrus.Fields{
				"uri":        v.URI,
				"status":     v.Status,
				"method":     v.Method,
				"remote_ip":  v.RemoteIP,
				"user_agent": v.UserAgent,
				"latency":    v.Latency.String(),
				"error":      v.Error,
			}).Info("request details")
			return nil
		},
	})
}
