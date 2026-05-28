package kafka

import (
	"accountService/kafkaEvents"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type NotificationProducer struct {
	notificationWriter *kafka.Writer
}

func NewProducer() *NotificationProducer {
	return &NotificationProducer{
		notificationWriter: &kafka.Writer{
			Addr:     kafka.TCP("kafka:9092"),
			Topic:    "notification",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *NotificationProducer) Close() error {
	if err := p.notificationWriter.Close(); err != nil {
		return err
	}
	return nil
}

func (p *NotificationProducer) SendSubscriptionCreated(ctx context.Context, event kafkaEvents.NotificationEvent) error {
	log.Println("Send Subscription Created Event")
	bytes, err := json.Marshal(event)
	if err != nil {
		log.Println("Marshal error:", err)
		return err
	}
	err = p.notificationWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.ProfileID.String()),
		Value: bytes,
	})
	if err != nil {
		log.Println("Write error:", err)
		return err
	}
	log.Println("Subscription Created Event sent successfully")
	return nil
}
