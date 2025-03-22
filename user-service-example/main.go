package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/zuyatna/edu-connect/user-service/database"
	"github.com/zuyatna/edu-connect/user-service/handler"
	"github.com/zuyatna/edu-connect/user-service/middlewares"
	"github.com/zuyatna/edu-connect/user-service/model"
	pb "github.com/zuyatna/edu-connect/user-service/pb/user"
	"github.com/zuyatna/edu-connect/user-service/repository"
	"github.com/zuyatna/edu-connect/user-service/routes"
	"github.com/zuyatna/edu-connect/user-service/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	db, err := gorm.Open(postgres.Open(database.InitDB()), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database!")
	}

	db.AutoMigrate(&model.User{})
	fmt.Println("database migrated")

	sigChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	quitChan := make(chan bool, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case err := <-errChan:
				logger.Fatalf("Server error: %v", err)
				quitChan <- true
			case <-sigChan:
				quitChan <- true
			}
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	grpcEndpoint := os.Getenv("GRPC_ENDPOINT")
	if grpcEndpoint == "" {
		grpcEndpoint = "localhost"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	go InitHTTPServer(errChan, port, grpcEndpoint, grpcPort)
	go InitGRPCServer(db, errChan, grpcEndpoint, grpcPort)

	<-quitChan
	logger.Info("Shutting down...")
}

func InitHTTPServer(errChan chan error, port, grpcEndpoint, grpcPort string) {
	conn, err := grpc.NewClient(grpcEndpoint+":"+grpcPort,
		grpc.WithInsecure(),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	userClient := pb.NewUserServiceClient(conn)

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${remote_ip} ${method} ${uri} ${status} ${latency_human}\n",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	userRoutes := routes.NewUserHTTPHandler(userClient)
	userRoutes.Routes(e)

	log.Info("Starting HTTP Server at port: ", port)
	errChan <- e.Start(":" + port)
}

func InitGRPCServer(db *gorm.DB, errChan chan error, grcpEndpoint, grpcPort string) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", grcpEndpoint, grpcPort))
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption
	if os.Getenv("ENV") == "production" {
		creds := credentials.NewServerTLSFromCert(&tls.Certificate{})
		opts = append(opts, grpc.Creds(creds))
	}

	opts = append(opts, grpc.UnaryInterceptor(middlewares.SelectiveAuthInterceptor))

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterUserServiceServer(grpcServer, userHandler)

	log.Info("Starting gRPC Server at", grcpEndpoint, ":", grpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		errChan <- err
	}
}
