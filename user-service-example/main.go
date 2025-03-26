package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"user-service-example/database"
	"user-service-example/handler"
	"user-service-example/middlewares"
	"user-service-example/model"
	pb "user-service-example/pb/user"
	"user-service-example/repository"
	"user-service-example/routes"
	"user-service-example/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
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
	var opts []grpc.DialOption

	if os.Getenv("ENV") == "production" {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(grpcEndpoint+":"+grpcPort, opts...)
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

func InitGRPCServer(db *gorm.DB, errChan chan error, grpcEndpoint, grpcPort string) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", grpcEndpoint, grpcPort))
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

	log.Info("Starting gRPC Server at", grpcEndpoint, ":", grpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		errChan <- err
	}
}
