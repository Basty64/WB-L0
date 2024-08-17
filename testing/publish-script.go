package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"wb/internal/models"
	"wb/internal/nats"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}
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
	orderData, err := getOrderDataFromFile("/Users/basty64/Programming/go/src/wb/testing/files/messages.json")
	if err != nil {
		log.Fatal(err)
	}

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

	// Публикация сообщения
	err = nats.PublishOrderToNATS(natsURL, natsSubject, data)
	if err != nil {
		log.Fatal(err)
	}
}

func getOrderDataFromFile(filename string) (models.Order, error) {
	msg, err := os.ReadFile(filename)
	if err != nil {
		return models.Order{}, err
	}

	var order models.Order
	err = json.Unmarshal(msg, &order)
	if err != nil {
		fmt.Println("Ошибка декодирования JSON:", err)
		return models.Order{}, err
	}
	return order, nil
}
