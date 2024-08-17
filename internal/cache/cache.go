package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"wb/internal/db/postgres"
	"wb/internal/models"
)

type Cache interface {
	InsertOrder(orderUid int, order models.Order) error
	GetOrder(orderUid int) (models.Order, error)
	GetAllOrders() ([]models.Order, error)
	LoadFromPostgres(ctx context.Context, database *postgres.RepoPostgres) error
}

type InMemoryCache struct {
	sync.Mutex
	Orders map[int]models.Order
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		Orders: make(map[int]models.Order),
	}
}

func (i *InMemoryCache) InsertOrder(orderUid int, order models.Order) error {
	i.Lock()
	defer i.Unlock()

	i.Orders[orderUid] = order

	return nil
}

var OrderNotFoundErr error

// GetOrder метод для веб-сервера - выдает заказ по айди
func (i *InMemoryCache) GetOrder(orderUid int) (models.Order, error) {
	i.Lock()
	defer i.Unlock()

	order, ok := i.Orders[orderUid]
	if !ok {
		OrderNotFoundErr = errors.New("order not found")
		return models.Order{}, fmt.Errorf("order not found: %w", OrderNotFoundErr)
	}

	return order, nil
}

// GetAllOrders показ всех заказов
func (i *InMemoryCache) GetAllOrders() ([]models.Order, error) {
	i.Lock()
	defer i.Unlock()

	var orders []models.Order
	for _, order := range i.Orders {
		orders = append(orders, order)
	}

	return orders, nil
}

// LoadFromPostgres загрузка всех заказов в кэш при запуске приложения
func (i *InMemoryCache) LoadFromPostgres(ctx context.Context, database *postgres.RepoPostgres) error {
	i.Lock()
	defer i.Unlock()

	orders, err := database.GetAllOrders(ctx)
	if err != nil {
		return fmt.Errorf("failed to load Orders from database: %w", err)
	}

	for _, order := range orders {
		i.Orders[order.OrderUid] = order
	}

	return nil
}
