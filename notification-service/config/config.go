package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var RabbitMQConn *amqp091.Connection

func InitDB() *gorm.DB {

	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	conStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)

	db, err := gorm.Open(postgres.Open(conStr), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db
	log.Println("Database connected successfully!")

	return db

}

func InitRabbitMQ() {

	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}

	RabbitMQUser := os.Getenv("MQUSER")
	RabbitMQPass := os.Getenv("MQPASS")
	RabbitMQHost := os.Getenv("MQHOST")
	RabbitMQPort := os.Getenv("MQPORT")
	RabbitMQVHost := os.Getenv("MQVHOST")

	conStr := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		RabbitMQUser, RabbitMQPass, RabbitMQHost, RabbitMQPort, RabbitMQVHost,
	)

	connection, err := amqp091.Dial(conStr)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	RabbitMQConn = connection
	fmt.Println("RabbitMQ connected")
}

func CloseRabbitMQ() {
	if RabbitMQConn != nil {
		RabbitMQConn.Close()
	}
}
