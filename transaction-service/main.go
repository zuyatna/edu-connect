package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"transaction-service/database"
	"transaction-service/handler"
	"transaction-service/middlewares"
	pbFuncCollect "transaction-service/pb/fund_collect"
	"transaction-service/pb/transaction"
	pbUser "transaction-service/pb/user"
	"transaction-service/repository"
	"transaction-service/routes"
	"transaction-service/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	ctx := context.Background()
	db := database.GetMongoDatabase()

	defer func() {
		if err := database.CloseMongoConnection(ctx); err != nil {
			logger.Fatalf("Failed to close MongoDB connection: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	quitChan := make(chan bool, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case err := <-errChan:
				logger.Fatalf("HTTP server error: %v", err)
				quitChan <- true
			case <-sigChan:
				logger.Info("Shutting down HTTP server")
				quitChan <- true
			}
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	grpcEndpoint := os.Getenv("GRPC_ENDPOINT")
	if grpcEndpoint == "" {
		grpcEndpoint = "localhost"
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50053"
	}

	transactionRepo := repository.NewTransactionRepository(db)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo)

	userConn, fundCollectConn := getServiceConnections()

	go InitHTTPServer(errChan, port, grpcEndpoint, grpcPort, transactionUsecase, userConn, fundCollectConn)
	go InitGRPCServer(db, errChan, grpcEndpoint, grpcPort, transactionUsecase, userConn, fundCollectConn)

	<-quitChan
	logger.Info("Shutting down...")

	userConn.Close()
	fundCollectConn.Close()
}

func getServiceConnections() (*grpc.ClientConn, *grpc.ClientConn) {
	grpcUserEndpoint := os.Getenv("GRPC_USER_ENDPOINT")
	if grpcUserEndpoint == "" {
		grpcUserEndpoint = "localhost"
	}
	grpcUserPort := os.Getenv("GRPC_USER_PORT")
	if grpcUserPort == "" {
		grpcUserPort = "50051"
	}

	// For Cloud Run, use TLS credentials
	var creds credentials.TransportCredentials
	if os.Getenv("ENV") == "production" {
		creds = credentials.NewClientTLSFromCert(nil, "")
	} else {
		creds = insecure.NewCredentials()
	}

	userConn, err := grpc.Dial(grpcUserEndpoint+":"+grpcUserPort,
		grpc.WithTransportCredentials(creds),
	)

	grpcInstitutionEndpoint := os.Getenv("GRPC_INSTITUTION_ENDPOINT")
	if grpcInstitutionEndpoint == "" {
		grpcInstitutionEndpoint = "localhost"
	}
	grpcInstitutionPort := os.Getenv("GRPC_INSTITUTION_PORT")
	if grpcInstitutionPort == "" {
		grpcInstitutionPort = "50052"
	}

	fundCollectConn, err := grpc.Dial(grpcInstitutionEndpoint+":"+grpcInstitutionPort,
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatalf("Failed to connect to fund collect service: %v", err)
	}

	return userConn, fundCollectConn
}

func InitHTTPServer(
	errChan chan error,
	port,
	grpcEndpoint,
	grpcPort string,
	transactionUsecase usecase.ITransactionUsecase,
	userConn *grpc.ClientConn,
	fundCollectConn *grpc.ClientConn,
) {
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

	transactionClient := transaction.NewTransactionServiceClient(conn)
	fundCollectClient := pbFuncCollect.NewFundCollectServiceClient(fundCollectConn)

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${remote_ip} ${method} ${uri} ${status} ${latency_human}\n",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	userClient := pbUser.NewUserServiceClient(userConn)
	paymentCallbackHandler := handler.NewPaymentCallbackHandler(transactionUsecase, userClient, fundCollectClient)
	e.POST("/api/payment/callback", func(c echo.Context) error {
		paymentCallbackHandler.HandleCallback(c.Response().Writer, c.Request())
		return nil
	})

	transactionRoutes := routes.NewTransactionHTTPHandler(transactionClient)
	transactionRoutes.Routes(e)

	log.Info("Starting HTTP Server at port: ", port)
	errChan <- e.Start(":" + port)
}

func InitGRPCServer(
	db *mongo.Database,
	errChan chan error,
	grpcEndpoint,
	grpcPort string,
	transactionUsecase usecase.ITransactionUsecase,
	userConn *grpc.ClientConn,
	fundCollectConn *grpc.ClientConn,
) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", grpcEndpoint, grpcPort))
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption
	if os.Getenv("ENV") == "production" {
		creds := credentials.NewServerTLSFromCert(&tls.Certificate{})
		opts = append(opts, grpc.Creds(creds))
	}

	opts = append(opts, grpc.UnaryInterceptor(middlewares.AuthGRPCInterceptor2))

	userClient := pbUser.NewUserServiceClient(userConn)
	fundCollectClient := pbFuncCollect.NewFundCollectServiceClient(fundCollectConn)

	transactionHandler := handler.NewTransactionHandler(transactionUsecase, userClient, fundCollectClient)

	transactionServer := grpc.NewServer(opts...)

	transaction.RegisterTransactionServiceServer(transactionServer, transactionHandler)

	log.Info("Starting gRPC Server at", grpcEndpoint, ":", grpcPort)
	if err := transactionServer.Serve(listener); err != nil {
		errChan <- err
	}
}
