package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"wb/internal/cache"
	"wb/internal/db/postgres"
	"wb/internal/handler"
	"wb/internal/nats"
)

type App struct {
	db         *postgres.RepoPostgres
	cache      *cache.InMemoryCache
	natsClient *nats.NatsStreamingClient
	server     *http.Server
	ctx        context.Context
}

func New() (*App, error) {

	return &App{}, nil

}

func (a *App) Run() error {

	ctx, cancel := context.WithCancel(context.Background())
	a.ctx = ctx
	defer cancel()

	// Получение параметров из окружения
	postgresURL := os.Getenv("DATABASE_URL_POSTGRES")
	natsClusterID := os.Getenv("NATS_CLUSTER_ID")
	natsClientID := os.Getenv("NATS_CLIENT_ID")
	natsURL := os.Getenv("NATS_URL")
	natsSubject := os.Getenv("NATS_SUBJECT")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	// Инициализация компонентов
	database, err := postgres.Connect(a.ctx, postgresURL)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	a.db = database
	defer a.db.CloseConnection()

	natsSubscriber, err := nats.NewNATSSubscriber(natsClusterID, natsClientID, natsURL)
	if err != nil {
		log.Fatalf("failed to initialize NATS subscriber: %v", err)
	}
	defer natsSubscriber.Close()

	// Инициализация кэша
	a.cache = cache.NewInMemoryCache()

	// Загрузка данных из базы данных в кэш при запуске
	if err := a.cache.LoadFromPostgres(a.ctx, a.db); err != nil {
		log.Printf("failed to load cache from database: %v", err)
	}

	// Подписка на канал в NATS Streaming
	if err := natsSubscriber.Subscribe(a.ctx, natsSubject, a.db, a.cache); err != nil {
		log.Fatalf("failed to subscribe to NATS: %v", err)
	}

	serverURL := fmt.Sprintf("%s:%s", host, port)
	server, err := handler.NewServer(a.ctx, serverURL, a.db, a.cache)

	go func() {
		if err := server.Start(serverURL); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start web server: %v", err)
		}
	}()

	<-ctx.Done()

	if err := server.Close(); err != nil {
		log.Fatalf("failed to shutdown web server: %v", err)
	}

	return nil
}
