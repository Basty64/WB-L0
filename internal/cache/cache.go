package cache

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"wb/internal/db/postgres"
	"wb/internal/models"
)

type Cache interface {
	InsertOrder(orderUid int, order models.Order) error
	GetOrder(orderUid int) (models.Order, bool)
	GetAllOrders() ([]models.Order, error)
	LoadFromPostgres(ctx context.Context, database *postgres.RepoPostgres) error
}

type InMemoryCache struct {
	mu     sync.Mutex
	Orders map[int]models.Order
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		Orders: make(map[int]models.Order),
	}
}

func (i *InMemoryCache) InsertOrder(orderUid int, order models.Order) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Паника в GetOrder: %v", r)
		}
	}()
	if i == nil {
		return errors.New("структура InMemoryCache не определена")
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	i.Orders[orderUid] = order

	return nil
}

var OrderNotFoundErr error

// GetOrder метод для веб-сервера - выдает заказ по айди
func (i *InMemoryCache) GetOrder(orderUid int) (models.Order, bool) {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Паника в GetOrder: %v", r)
		}
	}()
	if i == nil {
		return models.Order{}, false
	}

	i.mu.Lock()

	defer i.mu.Unlock()

	order, ok := i.Orders[orderUid]
	if !ok {
		OrderNotFoundErr = errors.New("order not found")
		return models.Order{}, false
	}

	return order, true
}

// GetAllOrders показ всех заказов
func (i *InMemoryCache) GetAllOrders() ([]models.Order, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Паника в GetOrder: %v", r)
		}
	}()

	var orders []models.Order
	for _, order := range i.Orders {
		orders = append(orders, order)
	}

	return orders, nil
}

// LoadFromPostgres загрузка всех заказов в кэш при запуске приложения
func (i *InMemoryCache) LoadFromPostgres(ctx context.Context, database *postgres.RepoPostgres) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Паника в GetOrder: %v", r)
		}
	}()

	orders, err := database.GetAllOrders(ctx)
	if err != nil {
		return fmt.Errorf(" Ошибка при загрузке заказов из хранилища: %w", err)
	}

	for _, order := range orders {
		i.Orders[order.ID] = order
	}

	return nil
}
