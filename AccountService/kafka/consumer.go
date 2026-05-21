package kafka

import (
	"accountService/kafkaEvents"
	"accountService/repository"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"strconv"
	"time"
)

const brokerAddress = "kafka:9092"

type Consumer struct {
	postCreatedReader *kafka.Reader
	postDeletedReader *kafka.Reader
}

func NewConsumer() *Consumer {
	return &Consumer{
		postCreatedReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{brokerAddress},
			Topic:   "post-created",
			GroupID: "account-post-created-group",

			MinBytes:              1,
			MaxBytes:              10e6,
			WatchPartitionChanges: true,
		},
		),
		postDeletedReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{brokerAddress},
			Topic:   "post-deleted",
			GroupID: "account-post-deleted-group",

			MinBytes:              1,
			MaxBytes:              10e6,
			WatchPartitionChanges: true,
		},
		),
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
			Topic:             "post-created",
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
		kafka.TopicConfig{
			Topic:             "post-deleted",
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	)
}

func (c *Consumer) ConsumePostCreated(r *repository.UsersRepository) {
	log.Println("Consuming Post Created Event started")
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := c.ensureTopics(ctx)
		cancel()
		if err == nil {
			break
		}
		log.Println("ensure kafka topics error:", err)
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
		log.Println("message received:", string(msg.Value))
		var event kafkaEvents.PostEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Println("unmarshal error:", err)
			if commitErr := c.postCreatedReader.CommitMessages(context.Background(), msg); commitErr != nil {
				log.Println("commit invalid message error:", commitErr)
			}
			continue
		}
		log.Println("event parsed:", event.EventID)
		ctx := context.Background()

		tx, err := r.DB.Begin(ctx)
		if err != nil {
			log.Println("begin tx error:", err)
			continue
		}
		processed, err := r.CheckProcessedEvent(ctx, tx, event.EventID)
		if err != nil {
			tx.Rollback(ctx)
			log.Println("check processed error:", err)
			continue
		}
		log.Println("processed:", processed)
		if processed {
			tx.Rollback(ctx)
			if err = c.postCreatedReader.CommitMessages(ctx, msg); err != nil {
				log.Println("commit processed event error:", err)
			}
			continue
		}
		err = r.IncrementPostsAmount(ctx, tx, event.ProfileID)
		if err != nil {
			tx.Rollback(ctx)
			log.Println("increment posts error:", err)
			continue
		}
		log.Println("posts incremented")
		err = r.SaveProcessedEvent(ctx, tx, event.EventID)
		if err != nil {
			tx.Rollback(ctx)
			log.Println("save processed event error:", err)
			continue
		}
		err = tx.Commit(ctx)
		if err != nil {
			log.Println("commit error:", err)
			continue
		}
		err = c.postCreatedReader.CommitMessages(ctx, msg)
		if err != nil {
			log.Println("commit kafka message error:", err)
			continue
		}
		log.Println("event handled successfully")
	}
}

func (c *Consumer) ConsumePostDeleted(r *repository.UsersRepository) {
	log.Println("Consuming Post Deleted Event started")
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		err := c.ensureTopics(ctx)
		cancel()
		if err == nil {
			break
		}
		log.Println("ensure kafka topics error:", err)
		time.Sleep(2 * time.Second)
	}
	for {
		log.Println("waiting for kafka message")
		msg, err := c.postDeletedReader.FetchMessage(context.Background())
		if err != nil {
			log.Println("read message error:", err)
			time.Sleep(time.Second)
			continue
		}
		log.Println("message received:", string(msg.Value))
		var event kafkaEvents.PostEvent
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Println("unmarshal error:", err)
			if commitErr := c.postDeletedReader.CommitMessages(context.Background(), msg); commitErr != nil {
				log.Println("commit invalid message error:", commitErr)
			}
			continue
		}
		log.Println("event parsed:", event.EventID)
		ctx := context.Background()
		tx, err := r.DB.Begin(ctx)
		if err != nil {
			log.Println("begin tx error:", err)
			continue
		}
		processed, err := r.CheckProcessedEvent(ctx, tx, event.EventID)
		if err != nil {
			tx.Rollback(ctx)
			log.Println("check processed error:", err)
			continue
		}
		log.Println("processed:", processed)
		if processed {
			tx.Rollback(ctx)
			if err = c.postDeletedReader.CommitMessages(ctx, msg); err != nil {
				log.Println("commit processed event error:", err)
			}
			continue
		}
		err = r.IncrementPostsAmount(ctx, tx, event.ProfileID)
		if err != nil {
			tx.Rollback(ctx)
			log.Println("decrement posts error:", err)
			continue
		}
		err = r.SaveProcessedEvent(ctx, tx, event.EventID)
		if err != nil {
			tx.Rollback(ctx)
			log.Println("save processed event error:", err)
			continue
		}
		err = tx.Commit(ctx)
		if err != nil {
			log.Println("commit error:", err)
			continue
		}
		err = c.postDeletedReader.CommitMessages(ctx, msg)
		if err != nil {
			log.Println("commit kafka message error:", err)
			continue
		}
		log.Println("event handled successfully")
	}
}
