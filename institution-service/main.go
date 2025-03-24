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
	"github.com/zuyatna/edu-connect/institution-service/database"
	"github.com/zuyatna/edu-connect/institution-service/handler"
	"github.com/zuyatna/edu-connect/institution-service/middlewares"
	"github.com/zuyatna/edu-connect/institution-service/model"
	"github.com/zuyatna/edu-connect/institution-service/pb/fund_collect"
	"github.com/zuyatna/edu-connect/institution-service/pb/institution"
	"github.com/zuyatna/edu-connect/institution-service/pb/post"
	"github.com/zuyatna/edu-connect/institution-service/repository"
	"github.com/zuyatna/edu-connect/institution-service/routes"
	"github.com/zuyatna/edu-connect/institution-service/usecase"
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

	initDB := database.GetDB()
	if initDB == nil {
		fmt.Println("Failed to initialize database")
		return
	}
	fmt.Println("Application started successfully")

	defer database.CloseDB()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: initDB}), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database!")
	}

	db.AutoMigrate(
		&model.Post{},
		&model.Institution{},
	)
	fmt.Println("Database migrated!")

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
		port = "8081"
	}

	grpcEndpoint := os.Getenv("GRPC_ENDPOINT")
	if grpcEndpoint == "" {
		grpcEndpoint = "localhost"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50061"
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

	insClient := institution.NewInstitutionServiceClient(conn)
	postClient := post.NewPostServiceClient(conn)

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${remote_ip} ${method} ${uri} ${status} ${latency_human}\n",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	insRoutes := routes.NewInstitutionHTTPHandler(insClient)
	insRoutes.Routes(e)

	postRoutes := routes.NewPostHTTPHandler(postClient)
	postRoutes.Routes(e)

	log.Info("Starting HTTP Server at port: ", port)
	errChan <- e.Start(":" + port)
}

func InitGRPCServer(db *gorm.DB, errChan chan error, grcpEndpoint, grpcPort string) {
	lis, err := net.Listen("tcp", grcpEndpoint+":"+grpcPort)
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption
	if os.Getenv("ENV") == "production" {
		creds := credentials.NewServerTLSFromCert(&tls.Certificate{})
		opts = append(opts, grpc.Creds(creds))
	}

	opts = append(opts, grpc.UnaryInterceptor(middlewares.SelectiveAuthInterceptor))

	insRepo := repository.NewInstitutionRepository(db)
	insUsecase := usecase.NewInstitutionUsecase(insRepo)
	insHandler := handler.NewInstitutionHandler(insUsecase)

	postRepo := repository.NewPostRepository(db)
	postUsecase := usecase.NewPostUsecase(postRepo)
	postHandler := handler.NewPostHandler(postUsecase)

	fundCollectRepo := repository.NewFundCollectRepository(db)
	fundCollectUsecase := usecase.NewFundCollectUsecase(fundCollectRepo)
	fundCollectHandler := handler.NewFundCollectHandler(fundCollectUsecase)

	grpcServer := grpc.NewServer(opts...)

	institution.RegisterInstitutionServiceServer(grpcServer, insHandler)
	post.RegisterPostServiceServer(grpcServer, postHandler)
	fund_collect.RegisterFundCollectServiceServer(grpcServer, fundCollectHandler)

	log.Info("Starting gRPC Server at", grcpEndpoint, ":", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		errChan <- err
	}
}
