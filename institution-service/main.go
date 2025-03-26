package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"institution-service/database"
	"institution-service/docs"
	"institution-service/handler"
	"institution-service/middlewares"
	"institution-service/model"
	"institution-service/pb/fund_collect"
	"institution-service/pb/institution"
	"institution-service/pb/post"
	"institution-service/repository"
	"institution-service/routes"
	"institution-service/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
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

	if err := db.AutoMigrate(&model.Institution{}); err != nil {
		logger.Fatalf("Failed to migrate Institution table: %v", err)
	}
	if err := db.AutoMigrate(&model.Post{}); err != nil {
		logger.Fatalf("Failed to migrate Post table: %v", err)
	}
	if err := db.AutoMigrate(&model.FundCollect{}); err != nil {
		logger.Fatalf("Failed to migrate FundCollect table: %v", err)
	}

	fmt.Println("Database migrated successfully!")

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
		grpcPort = "50052"
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
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.NewClient(grpcEndpoint+":"+grpcPort, opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	insClient := institution.NewInstitutionServiceClient(conn)
	postClient := post.NewPostServiceClient(conn)
	fundClient := fund_collect.NewFundCollectServiceClient(conn)

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${remote_ip} ${method} ${uri} ${status} ${latency_human}\n",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	docs.SwaggerInfo.Title = "EduConnect - Institution Service API Contract"
	docs.SwaggerInfo.Description = "This is a documentation EduConnect - Institution Service API Contract."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "institution-service-1011483964797.asia-southeast2.run.app"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"https"}
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	insRoutes := routes.NewInstitutionHTTPHandler(insClient)
	insRoutes.Routes(e)

	postRoutes := routes.NewPostHTTPHandler(postClient)
	postRoutes.Routes(e)

	fundCollectRoutes := routes.NewFundCollectHTTPHandler(fundClient)
	fundCollectRoutes.Routes(e)

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

	insRepo := repository.NewInstitutionRepository(db)
	insUsecase := usecase.NewInstitutionUsecase(insRepo)
	insHandler := handler.NewInstitutionHandler(insUsecase)

	postRepo := repository.NewPostRepository(db)
	postUsecase := usecase.NewPostUsecase(postRepo)
	postHandler := handler.NewPostHandler(postUsecase)

	fundCollectRepo := repository.NewFundCollectRepository(db)
	fundCollectUsecase := usecase.NewFundCollectUsecase(fundCollectRepo)
	fundCollectHandler := handler.NewFundCollectHandler(fundCollectUsecase, postUsecase)

	grpcServer := grpc.NewServer(opts...)

	institution.RegisterInstitutionServiceServer(grpcServer, insHandler)
	post.RegisterPostServiceServer(grpcServer, postHandler)
	fund_collect.RegisterFundCollectServiceServer(grpcServer, fundCollectHandler)

	log.Info("Starting gRPC Server at", grpcEndpoint, ":", grpcPort)
	if err := grpcServer.Serve(listener); err != nil {
		errChan <- err
	}
}
