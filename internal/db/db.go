package db

import (
	"context"
	"wb/internal/models"
)

type Database interface {
	InsertOrder(ctx context.Context, order *models.Order) (int, error)
	GetOrder(ctx context.Context, ID int) (models.Order, error)
	GetAllOrders(ctx context.Context) ([]models.Order, error)
	CloseConnection()
}
