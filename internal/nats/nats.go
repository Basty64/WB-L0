package nats

import (
	"context"
	"fmt"
	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
	"wb/internal/cache"
	"wb/internal/db/postgres"
	"wb/internal/models"
)

type NatsStreamingClient struct {
	conn  *stan.Conn
	cache *cache.InMemoryCache
	db    *postgres.RepoPostgres
	close chan bool
}

// NewNatsStreamingClient ...
func NewNatsStreamingClient(natsClusterID, natsURL, natsClientID string, db *postgres.RepoPostgres, oc *cache.InMemoryCache) (*NatsStreamingClient, error) {
	// Подключение к NATS Streaming
	sc, err := stan.Connect(natsClusterID, natsClientID, stan.NatsURL(natsURL))
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к NATS-Streaming: %w", err)
	}

	//NATS_SUBJECT = orders
	//NATS_CLIENT_ID = delivery_service

	return &NatsStreamingClient{
		conn:  &sc,
		cache: oc,
		db:    db,
		close: make(chan bool, 1),
	}, nil
}

func (nsc *NatsStreamingClient) Subscribe(ctx context.Context, natsSubject string) (stan.Subscription, error) {
	// Подписка на канал
	sub, err := (*nsc.conn).Subscribe(natsSubject,
		func(msg *stan.Msg) {
			log.Infoln("Received message from nats-streaming")
			order, err := models.NewOrder(msg.Data)
			if err != nil {
				log.Printf("ошибка при декодировании данных заказа: %v\n", err)
				return
			}

			id, err := nsc.db.InsertOrder(ctx, &order)
			if err.Error() == `ERROR: duplicate key value violates unique constraint "orders_pkey" (SQLSTATE 23505)` {
				log.Printf("Ошибка при записи данных заказа: %v\n", err)
				if err := msg.Ack(); err != nil {
					return
				}
			}
			if id == 0 {
				fmt.Println("Данные не записаны")
				if err := msg.Ack(); err != nil {
					return
				}
			} else {
				fmt.Printf("Заказ %d успешно записан в базу данных", id)
			}

			_, ok := nsc.cache.GetOrder(id)
			if ok {
				log.Printf("order with id: %d already in cache", id)

				if err := msg.Ack(); err != nil {
					return
				}
				return
			}
			err = nsc.cache.InsertOrder(order.ID, order)
			if err := msg.Ack(); err != nil {
				return
			}
			log.Printf("order with id: %s added in cache and database", id)

		}, stan.SetManualAckMode())
	if err != nil {
		return nil, err
	}
	return sub, nil
}

// Close ...
func (nsc *NatsStreamingClient) Close(sub stan.Subscription) error {
	nsc.close <- true

	conne := *nsc.conn

	if err := sub.Unsubscribe(); err != nil {
		return err
	}
	if err := conne.Close(); err != nil {
		return err
	}
	return nil
}
