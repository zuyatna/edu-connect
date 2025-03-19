package route

import (
	"os"
	"userService/handler"
	"userService/middleware"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func Init(e *echo.Echo, userHandler handler.UserHandler) {

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	v1 := e.Group("/v1")
	v1.Use(middleware.LogrusMiddleware(logger))
	user := v1.Group("/auth")
	user.POST("/register", userHandler.Register)
	user.POST("/login", userHandler.Login)

}
