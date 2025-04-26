// Updated main.go with Kafka integration
package main

import (
	"furniture-ecommerce/config"
	"furniture-ecommerce/kafka"
	"furniture-ecommerce/routes"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/stripe/stripe-go/v74"
)

func main() {
	
	config.InitDB()


	kafka.InitKafka()

	
	kafka.StartConsumers()

	
	r := routes.InitRoutes()


	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		stripeKey = "sk_test_51Q3VaJP8UEcGEhHsoMSaJV7OzABvuJDya5Vy0CHib5Z6V0Dee4MVVBAzx95NwOpMXQ5z0ylaCwC03xtVy2lteyfc00pGTq7MwR"
	}
	stripe.Key = stripeKey

	r.Static("/images", "./images")

	setupGracefulShutdown()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupGracefulShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")

		if kafka.Producer != nil {
			kafka.Producer.Close()
		}

		if config.DB != nil {
			sqlDB, err := config.DB.DB()
			if err == nil {
				sqlDB.Close()
			}
		}

		os.Exit(0)
	}()
}
