package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"wb/internal/models"
	"wb/internal/nats"
)

func main() {

	natsClusterID := os.Getenv("NATS_CLUSTER_ID")
	natsClientID := os.Getenv("NATS_CLIENT_ID")
	natsURL := os.Getenv("NATS_URL")
	natsSubject := os.Getenv("NATS_SUBJECT")

	// Подключение к серверу NATS Streaming
	nc, err := stan.Connect(natsClusterID, natsClientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatal(err)
	}
	defer func(nc stan.Conn) {
		err := nc.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(nc)

	// Чтение данных из файла
	orderData := getOrderDataFromFile("order.json")

	// Создание сообщения
	message := models.Order{
		OrderUid:          orderData.OrderUid,
		TrackNumber:       orderData.TrackNumber,
		Entry:             orderData.Entry,
		Delivery:          orderData.Delivery,
		Payment:           orderData.Payment,
		Item:              orderData.Item,
		Locale:            orderData.Locale,
		InternalSignature: orderData.InternalSignature,
		CustomerId:        orderData.CustomerId,
		DeliveryService:   orderData.DeliveryService,
		Shardkey:          orderData.Shardkey,
		SmId:              orderData.SmId,
		DateCreated:       orderData.DateCreated,
		OofShard:          orderData.OofShard,
	}

	// Преобразование сообщения в JSON
	data, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	publishOrder(natsURL, natsSubject, orderData)

	// Публикация сообщения
	err = nats.PublishOrderToNATS(natsURL, natsSubject, data)
	if err != nil {
		log.Fatal(err)
	}
}

func getOrderDataFromFile(filename string) models.Order {
	return models.Order{}
}

func publishOrder(natsStreamingURL string, subject string, orderData models.Order) {
	// ... (реализация публикации данных в канал)
}
