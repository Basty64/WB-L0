package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"time"
	"wb/internal/models"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	// Генерируем уникальный идентификатор клиента (например, используя UUID)
	clientID := "client-" + fmt.Sprintf("%d", time.Now().UnixNano())

	natsURL := os.Getenv("NATS_URL")
	natsSubject := os.Getenv("NATS_SUBJECT")

	// Подключение к серверу NATS Streaming
	nc, err := stan.Connect("test-cluster", clientID, stan.NatsURL(natsURL))
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
	orderData, err := getOrderDataFromFile("./testing/files/messages.json")
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
	if err := nc.Publish(natsSubject, data); err != nil {
		fmt.Errorf("failed to publish order to NATS: %w", err)
	}

	log.Printf("published order to subject: %s", natsSubject)
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
