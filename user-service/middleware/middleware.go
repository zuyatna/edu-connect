package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

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

func ValidateJWT(tokenString string) (string, error) {

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	JWTSecret := os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Error("JWT_SECRET is not set")
		return "", fmt.Errorf("server error: JWT secret is missing")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Warn("Unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(JWTSecret), nil
	})

	if err != nil {
		log.WithFields(logrus.Fields{
			"token": tokenString,
			"error": err.Error(),
		}).Error("JWT parsing failed")
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			log.Warn("Email not found in token claims")
			return "", fmt.Errorf("email not found in token")
		}

		log.WithFields(logrus.Fields{
			"email": email,
		}).Info("JWT validated successfully")
		return email, nil
	}

	log.Warn("Invalid token")
	return "", fmt.Errorf("invalid token")
}
