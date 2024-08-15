package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
	"log"
	"os"
	"time"
	"wb/internal/db"
	"wb/internal/models"
)

type NatsStreamingClient struct {
	conn  stan.Conn
	sub   stan.Subscription
	close chan bool
}

func NewNatsStreamingClient(ctx context.Context, subject string, clusterID string, clientID string, url string, db db.Database) (*NatsStreamingClient, error) {

	// Подключение к NATS-Streaming
	conn, err := stan.Connect(clusterID, clientID, stan.NatsURL(url))
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к NATS-Streaming: %w", err)
	}

	// Подписка на канал
	sub, err := conn.Subscribe(subject, func(m *stan.Msg) {
		var message models.Order

		err := json.Unmarshal(m.Data, &message)
		if err != nil {
			log.Fatal(err)
		}

		// Обработка данных заказа
		// TO DO:
		// Добавить валидацию данных заказа, чтобы нельзя было кидать в канал что угодно

		// Сохранение данных в базу данных
		// ...

		_, err = db.InsertOrder(ctx, message)
		if err != nil {
			fmt.Errorf("Ошибка при сохранении сообщения из nats-streaming: %d", err)
		}

	}, stan.StartAt(pb.StartPosition_NewOnly))
	if err != nil {
		return nil, fmt.Errorf("Ошибка при подписке на канал: %w", err)
	}

	return &NatsStreamingClient{
		conn:  conn,
		sub:   sub,
		close: make(chan bool, 1),
	}, nil
}

func (c *NatsStreamingClient) Close() error {
	c.close <- true
	c.sub.Unsubscribe()
	if err := c.conn.Close(); err != nil {
		return err
	}
	return nil
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func PublishOrderToNATS(natsURL string, subject string, order []byte) error {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return fmt.Errorf("Ошибка при подключении к  NATS: %w", err)
	}
	defer nc.Close()

	if err := nc.Publish(subject, order); err != nil {
		return fmt.Errorf("failed to publish order to NATS: %w", err)
	}

	log.Printf("published order to subject: %s", subject)

	return nil
}
