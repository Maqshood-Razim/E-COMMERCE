package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)


var Producer *kafka.Writer


func InitKafka() {
	Producer = &kafka.Writer{
		Addr:     kafka.TCP("localhost:29092"),
		Balancer: &kafka.LeastBytes{},
	}

	fmt.Println("Kafka producer initialized")
}


func PublishEvent(topic string, event interface{}) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = Producer.WriteMessages(ctx,
		kafka.Message{
			Topic: topic,
			Value: eventJSON,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish event to %s: %w", topic, err)
	}

	return nil
}
