package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"frauddetector/internal/api"
	"frauddetector/internal/storage"
	"frauddetector/internal/ml"
	"frauddetector/internal/kafka"
)

func main() {
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
    kafkaCfg := kafka.Config{
        Brokers:     brokers,
        GroupID:     os.Getenv("KAFKA_GROUP_ID"),
        Topic:       os.Getenv("KAFKA_TOPIC"),
        AlertsTopic: os.Getenv("KAFKA_ALERTS_TOPIC"),
    }
	//initializing Redis client
	redisAddr := os.Getenv("REDIS_ADDR")
	redisClient := storage.NewRedisClient(redisAddr)

	// ML client
    mlClient := ml.NewClient(os.Getenv("ML_URL"))

	kafkaCfg := kafka.Config{
		Brokers: []string{"localhost:9092"},
		Topic:   "transactions",
		GroupID: "fraud-detector-group",
		AlertsTopic: "fraud-alerts",
	}
	consumer := kafka.NewConsumer(kafkaCfg, redisClient, mlClient)
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    go func() {
        if err := consumer.Run(ctx); err != nil {
            log.Fatalf("Kafka consumer error: %v", err)
        }
    }()


	//create Gin router
	router := api.NewRouter(redisClient, mlClient)

	// Handle graceful shutdown on SIGINT/SIGTERM
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-quit
        log.Println("Shutting down consumer...")
        cancel()
    }()

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}