package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
	"wb/internal/models"
)

func main() {

	err := godotenv.Load("/wb/.env")
	if err != nil {
		log.Error("Error loading .env file", err)
	}

	// Генерируем уникальный идентификатор клиента
	clientID := "client-" + fmt.Sprintf("%d", time.Now().UnixNano())

	natsURL := os.Getenv("NATS_URL")
	natsSubject := os.Getenv("NATS_SUBJECT")

	path := "/wb/testing/files/messages"

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

	messages, err := os.ReadDir(path)
	if err != nil {
		log.Error(err)
	}

	for _, m := range messages {

		workpath := filepath.Join(path, m.Name())

		// Чтение данных из файла
		orderData, err := getOrderDataFromFile(workpath)
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
			Items:             orderData.Items,
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
			log.Errorf("failed to publish order to NATS: %s", err)
		}

		log.Printf("published order to subject: %s, name: %s", natsSubject, m.Name())
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
		log.Error("Ошибка декодирования JSON:", err)
		return models.Order{}, err
	}
	return order, nil
}
