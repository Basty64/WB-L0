package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func PublishOrderToNATS(natsURL string, subject string, order []byte) error {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}
	defer nc.Close()

	if err := nc.Publish(subject, order); err != nil {
		return fmt.Errorf("failed to publish order to NATS: %w", err)
	}

	log.Printf("published order to subject: %s", subject)

	return nil
}

func SendOrderToHTTPServer(url string, order []byte) error {
	req, err := http.NewRequest(http.MethodPost, url, ioutil.NopCloser(order))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send order to HTTP server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("HTTP request failed with status code: %d, body: %s", resp.StatusCode, string(body))
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
