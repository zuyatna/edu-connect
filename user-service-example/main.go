package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/zuyatna/edu-connect/user-service/database"
	"go.mongodb.org/mongo-driver/mongo"
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
	
}

func InitGRPCServer(db *mongo.Database, errChan chan error) {

}
