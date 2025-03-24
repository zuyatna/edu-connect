package database

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	database   *mongo.Database
	client     *mongo.Client
	clientLock sync.Mutex
	once       sync.Once
	log        = logrus.New()
)

func GetMongoClient() *mongo.Client {
	clientLock.Lock()
	defer clientLock.Unlock()

	if client != nil {
		return client
	}

	once.Do(func() {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})

		err := godotenv.Load()
		if err != nil {
			log.Fatal("Failed to load .env file")
		}

		uri := os.Getenv("MONGO_URI")
		if uri == "" {
			log.Fatal("MONGO_URI is not set")
		}

		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().
			ApplyURI(uri).
			SetServerAPIOptions(serverAPI).
			SetMaxPoolSize(10).
			SetMinPoolSize(4).
			SetMaxConnIdleTime(60 * time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err = mongo.Connect(ctx, opts)
		if err != nil {
			log.Fatalf("Failed to connect to MongoDB: %v", err)
		}

		if err := client.Ping(context.Background(), nil); err != nil {
			log.Fatalf("Failed to ping MongoDB: %v", err)
		}

		log.Info("Connected to MongoDB")
	})

	return client
}

func GetMongoDatabase() *mongo.Database {
	if database != nil {
		return database
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env file")
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		log.Fatal("MONGO_DB is not set")
	}

	database = GetMongoClient().Database(dbName)

	return database
}

func CloseMongoConnection(ctx context.Context) error {
	clientLock.Lock()
	defer clientLock.Unlock()

	if client != nil {
		err := client.Disconnect(ctx)
		if err != nil {
			return err
		}
		client = nil
		database = nil
		log.Info("MongoDB connection closed")
	}

	return nil
}
