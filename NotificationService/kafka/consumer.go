package kafka

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"notificationService/kafkaEvents"
	"notificationService/models"
	"notificationService/repository"
	"strconv"
	"time"
)

const brokerAddress = "kafka:9092"

type Consumer struct {
	notificationReader *kafka.Reader
}

func NewConsumer() *Consumer {
	return &Consumer{
		notificationReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{brokerAddress},
			Topic:   "notification",
			GroupID: "notification-group",

			MinBytes:              10e3,
			MaxBytes:              10e6,
			WatchPartitionChanges: true,
		}),
	}
}

func (c *Consumer) ensureTopics(ctx context.Context) error {
	conn, err := kafka.DialContext(ctx, "tcp", brokerAddress)
	if err != nil {
		return err
	}
	defer conn.Close()
	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	controllerConn, err := kafka.DialContext(ctx, "tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()
	return controllerConn.CreateTopics(
		kafka.TopicConfig{
			Topic:             "notification",
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	)
}

func (c *Consumer) ConsumeNotification(r *repository.NotificationRepository) error {
	log.Println("Consuming Notification Event started")
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := c.ensureTopics(ctx)
		cancel()
		if err == nil {
			break
		}
		log.Println("Ensure kafka topics error:", err)
		time.Sleep(2 * time.Second)
	}
	for {
		log.Println("waiting for kafka message")
		msg, err := c.notificationReader.FetchMessage(context.Background())
		if err != nil {
			log.Println("read message error:", err)
			time.Sleep(time.Second)
			continue
		}
		log.Println("message received:", string(msg.Value))
		var event kafkaEvents.NotificationEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Println("unmarshal message error:", err)
			if commitErr := c.notificationReader.CommitMessages(context.Background(), msg); commitErr != nil {
				log.Println("commit invalid message error:", commitErr)
			}
			continue
		}
		log.Println("event parsed:", event.EventID)
		ctx := context.Background()
		tx, err := r.DB.Begin(ctx)
		if err != nil {
			log.Println("start transaction error:", err)
			continue
		}
		processed, err := r.CheckEventProcessed(ctx, tx, event.EventID)
		if err != nil {
			tx.Rollback(ctx)
			log.Println("check processed error:", err)
			continue
		}
		log.Println("processed:", processed)
		if processed {
			tx.Rollback(ctx)
			if err := c.notificationReader.CommitMessages(ctx, msg); err != nil {
				log.Println("commit processed event error:", err)
			}
			continue
		}
		not := models.Notification{
			NotificationID: uuid.New(),
			ProfileID:      event.ProfileID,
			ActorID:        event.ActorID,
			Type:           event.Type,
			EntityID:       event.EntityID,
			CreatedAt:      event.CreatedAt,
		}
		err = r.SaveNotification(ctx, tx, not)
		if err != nil {
			tx.Rollback(ctx)
			log.Println("save notification error:", err)
			continue
		}
		err = r.ProcessEvent(ctx, tx, event.EventID)
		if err != nil {
			tx.Rollback(ctx)
			log.Println("save processed event error:", err)
			continue
		}
		err = tx.Commit(ctx)
		if err != nil {
			log.Println("commit transaction error:", err)
			continue
		}
		err = c.notificationReader.CommitMessages(ctx, msg)
		if err != nil {
			log.Println("commit kafka msg error:", err)
			continue
		}
		log.Println("kafka event handled successfully")
	}
}
