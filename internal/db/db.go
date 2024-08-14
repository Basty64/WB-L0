package db

import (
	"context"
	"wb/internal/models"
)

type Database interface {
	Connect(ctx context.Context, url string) error
	InsertOrder(ctx context.Context, order models.Order) (int, error)
	GetOrder(ctx context.Context, orderUid int) (models.Order, error)
	GetAllOrders(ctx context.Context) ([]models.Order, error)
	CloseConnection()
}
