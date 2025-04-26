package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func StartConsumers() {
	go consumeCartEvents()
	go consumeOrderEvents()
	go consumePaymentEvents()

	log.Println("Kafka consumers started")
}


func consumeCartEvents() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:29092"},
		Topic:    "cart_events",
		GroupID:  "cart-service",
		MinBytes: 10e3, 
		MaxBytes: 10e6, 
		MaxWait:  time.Second,
	})
	defer reader.Close()

	ctx := context.Background()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Cart consumer error: %v\n", err)
			continue
		}

		var event CartEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to parse cart event: %v\n", err)
			continue
		}

		log.Printf("Processing cart event: User %d %s product %d (qty: %d, price: %d)\n",
			event.UserID, event.Action, event.ProductID, event.Quantity, event.Price)
	}
}

func consumeOrderEvents() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:29092"},
		Topic:    "order_events",
		GroupID:  "order-service",
		MinBytes: 10e3, 
		MaxBytes: 10e6, 
		MaxWait:  time.Second,
	})
	defer reader.Close()

	ctx := context.Background()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Order consumer error: %v\n", err)
			continue
		}

		var event OrderEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to parse order event: %v\n", err)
			continue
		}

		log.Printf("Processing order event: Order %d for user %d (status: %s, amount: %.2f)\n",
			event.OrderID, event.UserID, event.Status, event.TotalPrice)
	}
}


func consumePaymentEvents() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:29092"},
		Topic:    "payment_events",
		GroupID:  "payment-service",
		MinBytes: 10e3, 
		MaxBytes: 10e6, 
		MaxWait:  time.Second,
	})
	defer reader.Close()

	ctx := context.Background()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Payment consumer error: %v\n", err)
			continue
		}

		var event PaymentEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to parse payment event: %v\n", err)
			continue
		}

		log.Printf("Processing payment event: Order %d payment %s (method: %s, amount: %.2f)\n",
			event.OrderID, event.Status, event.PaymentMethod, event.Amount)
	}
}
