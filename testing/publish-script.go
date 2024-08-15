package main

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"log"
	"wb/internal/models"
)

func main() {
	natsStreamingURL := "nats://localhost:4223"
	subject := "orders"

	// Подключение к серверу NATS Streaming
	nc, err := stan.Connect("my_cluster_id", "my_client_id", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

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
	publishOrder(natsStreamingURL, subject, orderData)

	// Публикация сообщения
	err = nc.Publish("my_subject", data)
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
