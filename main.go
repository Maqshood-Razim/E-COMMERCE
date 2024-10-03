package main

import (
	"furniture-ecommerce/config"
	"furniture-ecommerce/routes"
	"os"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

func main() {
	// Load environment variables
	config.InitDB()
	r := routes.InitRoutes()

	stripe.Key = os.Getenv("sk_test_51Q3VaJP8UEcGEhHsoMSaJV7OzABvuJDya5Vy0CHib5Z6V0Dee4MVVBAzx95NwOpMXQ5z0ylaCwC03xtVy2lteyfc00pGTq7MwRZ") // or directly assign your key here

	// Optionally, you can create a client
	sc := &client.API{}
	sc.Init(stripe.Key, nil)

	r.Static("/images", "./images")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
