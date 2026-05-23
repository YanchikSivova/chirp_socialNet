package kafka

import (
	"context"
	"encoding/json"
	"feedService/kafkaEvents"
	"feedService/service"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"strconv"
	"time"
)

const brokerAddress = "kafka:9092"

type Consumer struct {
	postCreatedReader *kafka.Reader
	service           *service.FeedService
}

func NewConsumer(s *service.FeedService) *Consumer {
	return &Consumer{postCreatedReader: kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   "post-created",
		GroupID: "feed-post-created-group",

		MinBytes:              1,
		MaxBytes:              10e6,
		WatchPartitionChanges: true,
	}),
		service: s}
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
			Topic:             "post_created",
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
}

func (c *Consumer) ConsumePostCreated() {
	log.Println("Consuming post created event started")
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := c.ensureTopics(ctx)
		cancel()
		if err == nil {
			break
		}
		log.Println("ensure kafka topic error:", err)
		time.Sleep(2 * time.Second)
	}
	for {
		log.Println("waiting for kafka message")
		msg, err := c.postCreatedReader.FetchMessage(context.Background())
		if err != nil {
			log.Println("read message error:", err)
			time.Sleep(time.Second)
			continue
		}
		log.Println("kafka message:", string(msg.Value))
		var event kafkaEvents.PostEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Println("unmarshal error:", err)
			continue
		}
		log.Println("event parsed:", event.EventID)
		ctx := context.Background()
		err = c.service.ProcessNewPost(event)
		if err != nil {
			log.Println("process post error:", err)
			continue
		}
		err = c.postCreatedReader.CommitMessages(ctx, msg)
		if err != nil {
			log.Println("commit messages error:", err)
			continue
		}
		log.Println("commit messages processed")
	}
}
