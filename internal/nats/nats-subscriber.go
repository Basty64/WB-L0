package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
	"log"
	"wb/internal/cache"
	"wb/internal/db"
	"wb/internal/models"
)

type NatsSubscriber interface {
	Subscribe(ctx context.Context, subject string, database db.Database, cache cache.Cache) error
}

type Subscriber struct {
	conn stan.Conn
	nc   *nats.Conn
}

func NewNATSSubscriber(clusterID, clientID, natsURL string) (*Subscriber, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	conn, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS Streaming: %w", err)
	}

	return &Subscriber{
		conn: conn,
		nc:   nc,
	}, nil
}

func (n *Subscriber) Subscribe(ctx context.Context, subject string, database db.Database, cache *cache.InMemoryCache) error {
	sub, err := n.conn.Subscribe(subject, func(msg *stan.Msg) {
		var order models.Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("failed to unmarshal order from NATS message: %v", err)
			return
		}

		if _, err := database.InsertOrder(context.Background(), order); err != nil {
			log.Printf("failed to insert order into database: %v", err)
			return
		}

		if err := cache.InsertOrder(order.OrderUid, order); err != nil {
			log.Printf("failed to set order in cache: %v", err)
			return
		}

		log.Printf("received order with order_uid: %s", order.OrderUid)
	}, stan.DurableName("durable-order-subscription"), stan.SetManualAckMode())
	if err != nil {
		return fmt.Errorf("failed to subscribe to subject: %w", err)
	}

	defer sub.Close()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				log.Printf("failed to ack message: %v", err)

			}
		}
	}()

	return nil
}

func (n *Subscriber) Close() error {
	if err := n.conn.Close(); err != nil {
		return fmt.Errorf("failed to close NATS Streaming connection: %w", err)
	}

	//FIX THE ERR, it must be variable
	if err := "close nc"; err != "" {
		n.nc.Close()
		return fmt.Errorf("failed to close NATS connection: %w", err)
	}

	return nil
}
