package nats

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
	"log"
	"wb/internal/db/postgres"
	"wb/internal/models"
)

type NatsStreamingClient struct {
	conn  stan.Conn
	sub   stan.Subscription
	close chan bool
}

// NewNatsStreamingClient ...
func NewNatsStreamingClient(ctx context.Context, natsClusterID, natsURL, natsSubject string, database *postgres.RepoPostgres) (*NatsStreamingClient, error) {
	// Подключение к NATS Streaming
	conn, err := stan.Connect(natsClusterID, natsSubject, stan.NatsURL(natsURL))
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к NATS-Streaming: %w", err)
	}

	// Подписка на канал
	sub, err := conn.Subscribe(natsSubject, func(msg *stan.Msg) {
		order, err := models.NewOrder(msg.Data)
		if err != nil {
			log.Printf("ошибка при декодировании данных заказа: %v\n", err)
			return
		}

		id, err := database.InsertOrder(ctx, order)
		if err != nil {
			log.Printf("Ошибка при записи данных заказа: %v\n", err)
		}
		if id == 0 {
			fmt.Println("Данные не записаны")
		} else {
			fmt.Printf("Заказ %d успешно записан", id)
		}
	}, stan.StartAt(pb.StartPosition_NewOnly))
	if err != nil {
		return nil, fmt.Errorf("ошибка при подписке на канал: %w", err)
	}

	return &NatsStreamingClient{
		conn:  conn,
		sub:   sub,
		close: make(chan bool, 1),
	}, nil
}

// Close ...
func (c *NatsStreamingClient) Close() error {
	c.close <- true
	if err := c.sub.Unsubscribe(); err != nil {
		return err
	}
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}

func PublishOrderToNATS(natsURL string, subject string, order []byte) error {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return fmt.Errorf("ошибка при подключении к NATS: %w", err)
	}
	defer nc.Close()

	if err := nc.Publish(subject, order); err != nil {
		return fmt.Errorf("failed to publish order to NATS: %w", err)
	}

	log.Printf("published order to subject: %s", subject)

	return nil
}
