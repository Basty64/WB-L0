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
	SetOrder(ctx context.Context, orderUid int, order models.Order) error
	GetOrder(ctx context.Context, orderUid int) (*models.Order, error)
	GetAllOrders(ctx context.Context) ([]models.Order, error)
	LoadFromPostgres(ctx context.Context, database postgres.RepoPostgres) error
}

type InMemoryCache struct {
	sync.Mutex
	orders map[int]models.Order
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		orders: make(map[int]models.Order),
	}
}

func (c *InMemoryCache) SetOrder(ctx context.Context, orderUid int, order models.Order) error {
	c.Lock()
	defer c.Unlock()

	c.orders[orderUid] = order

	return nil
}

func (c *InMemoryCache) GetOrder(ctx context.Context, orderUid int) (*models.Order, error) {
	c.Lock()
	defer c.Unlock()

	order, ok := c.orders[orderUid]
	if !ok {
		return nil, fmt.Errorf("order not found: %w", errors.New("order not found"))
	}

	return &order, nil
}

func (c *InMemoryCache) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	c.Lock()
	defer c.Unlock()

	var orders []models.Order
	for _, order := range c.orders {
		orders = append(orders, order)
	}

	return orders, nil
}

func (c *InMemoryCache) LoadFromPostgres(ctx context.Context, database *postgres.RepoPostgres) error {
	c.Lock()
	defer c.Unlock()

	orders, err := database.GetAllOrders(ctx)
	if err != nil {
		return fmt.Errorf("failed to load orders from database: %w", err)
	}

	for _, order := range orders {
		c.orders[order.OrderUid] = order
	}

	return nil
}
