package kafka

import (
	"context"
	"encoding/json"
	kafkaGo "github.com/segmentio/kafka-go"
	"log"
	"postService/kafkaEvents"
)

type Producer struct {
	postCreatedWriter  *kafkaGo.Writer
	postDeletedWriter  *kafkaGo.Writer
	notificationWriter *kafkaGo.Writer
}

func NewProducer() *Producer {
	return &Producer{
		postCreatedWriter: &kafkaGo.Writer{
			Addr:     kafkaGo.TCP("kafka:9092"),
			Topic:    "post-created",
			Balancer: &kafkaGo.LeastBytes{},
		},
		postDeletedWriter: &kafkaGo.Writer{
			Addr:     kafkaGo.TCP("kafka:9092"),
			Topic:    "post-deleted",
			Balancer: &kafkaGo.LeastBytes{},
		},
		notificationWriter: &kafkaGo.Writer{
			Addr:     kafkaGo.TCP("kafka:9092"),
			Topic:    "notification",
			Balancer: &kafkaGo.LeastBytes{},
		},
	}
}

func (p *Producer) Close() error {
	if err := p.postCreatedWriter.Close(); err != nil {
		return err
	}
	if err := p.postDeletedWriter.Close(); err != nil {
		return err
	}
	if err := p.notificationWriter.Close(); err != nil {
		return err
	}
	return nil
}

func (p *Producer) SendPostCreated(ctx context.Context, event kafkaEvents.PostEvent) error {
	log.Println("Sending post created event")
	bytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Marshal error: %v", err)
		return err
	}
	//
	err = p.postCreatedWriter.WriteMessages(ctx, kafkaGo.Message{
		Key:   []byte(event.ProfileID.String()),
		Value: bytes,
	})

	if err != nil {
		log.Printf("WriteMessages error: %v", err)
		return err
	}

	log.Println("Kafka message sent successfully")

	return nil
}

func (p *Producer) SendPostDeleted(ctx context.Context, event kafkaEvents.PostEvent) error {
	log.Println("Deleting post created event")
	bytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Marshal error: %v", err)
		return err
	}
	err = p.postDeletedWriter.WriteMessages(ctx, kafkaGo.Message{
		Key:   []byte(event.ProfileID.String()),
		Value: bytes,
	})
	if err != nil {
		log.Printf("WriteMessages error: %v", err)
		return err
	}
	log.Println("Kafka message sent successfully")
	return nil
}

func (p *Producer) SendNotification(ctx context.Context, event kafkaEvents.NotificationEvent) error {
	log.Println("Send Subscription Created Event")
	bytes, err := json.Marshal(event)
	if err != nil {
		log.Println("Marshal error:", err)
		return err
	}
	err = p.notificationWriter.WriteMessages(ctx, kafkaGo.Message{
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
