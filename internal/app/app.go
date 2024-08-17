package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
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

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	// Получение параметров из окружения
	postgresURL := os.Getenv("DATABASE_URL_POSTGRES")
	natsClusterID := os.Getenv("NATS_CLUSTER_ID")
	//natsClientID := os.Getenv("NATS_CLIENT_ID")
	natsURL := os.Getenv("NATS_URL")
	natsSubject := os.Getenv("NATS_SUBJECT")
	host := os.Getenv("APP_HOST")
	port := os.Getenv("APP_PORT")

	// Инициализация компонентов

	// База данных
	database, err := postgres.Connect(a.ctx, postgresURL)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	a.db = database
	defer a.db.CloseConnection()

	// Обработка сообщений с помощью Nats-streaming
	natsClient, err := nats.NewNatsStreamingClient(a.ctx, natsClusterID, natsURL, natsSubject, a.db)
	if err != nil {
		log.Fatalf("failed to initialize NATS subscriber: %v", err)
	}

	defer func(natsSubscriber *nats.NatsStreamingClient) {
		err := natsSubscriber.Close()
		if err != nil {
		}
	}(natsClient)

	// Обработка сигнала прерывания
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	//<-c

	// Инициализация кэша
	a.cache = cache.NewInMemoryCache()

	// Загрузка данных из базы данных в кэш при запуске
	if err := a.cache.LoadFromPostgres(a.ctx, a.db); err != nil {
		log.Printf("failed to load cache from database: %v", err)
	}

	serverURL := fmt.Sprintf("%s:%s", host, port)
	server, err := handler.NewServer(a.ctx, serverURL, a.db, a.cache)

	go func() {
		if err := server.Start(serverURL); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "Failed to start HTTP server: %v\n", err)
			os.Exit(1)
		}

	}()

	<-ctx.Done()

	if err := server.Close(); err != nil {
		log.Fatalf("failed to shutdown web server: %v", err)
	}

	return nil
}
