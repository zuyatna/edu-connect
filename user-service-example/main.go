package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"github.com/zuyatna/edu-connect/user-service/database"
	"github.com/zuyatna/edu-connect/user-service/model"
	pb "github.com/zuyatna/edu-connect/user-service/pb/user"
	"github.com/zuyatna/edu-connect/user-service/routes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	go InitHTTPServer(errChan, port)
	go InitGRPCServer(db, errChan)

	<-quitChan
	logger.Info("Shutting down...")
}

func InitHTTPServer(errChan chan error, port string) {
	grpcEndpoint := os.Getenv("GRPC_ENDPOINT")
	if grpcEndpoint == "" {
		grpcEndpoint = "localhost:50051"
	}

	conn, err := grpc.NewClient(":"+grpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close gRPC connection: %v", err)
		}
	}(conn)

	userClient := pb.NewUserServiceClient(conn)

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${remote_ip} ${method} ${uri} ${status} ${latency_human}\n",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// TODO: refactor databases to use postgres
	userRoutes := routes.NewUserHTTPHandler(userClient)
	userRoutes.Routes(e)

	log.Info("Starting HTTP Server at port: ", port)
	errChan <- e.Start(":" + port)
}

// TODO: refactor databases to use postgres
func InitGRPCServer(db *gorm.DB, errChan chan error) {

}
