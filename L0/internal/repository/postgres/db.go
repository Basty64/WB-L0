package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"wb/L0/internal/models"
)

type RepoPostgres struct {
	connection *pgxpool.Pool
}

type Repository interface {
	InsertOrder(ctx context.Context, order models.Order) error
	GetOrder(ctx context.Context, orderUid string) (*models.Order, error)
	GetAllOrders(ctx context.Context) ([]models.Order, error)
}

func NewRepository(conn *pgxpool.Pool) *RepoPostgres {
	return &RepoPostgres{connection: conn}
}

func (repo *RepoPostgres) Close() {
	repo.connection.Close()
}

func (repo *RepoPostgres) InsertOrder(ctx context.Context, order models.Order) (int, error) {

	err := repo.connection.QueryRow(ctx, "INSERT INTO orders (TrackNumber,Entry,Delivery,Payment,Cart,Locale,InternalSignature,CustomerId,DeliveryService,Shardkey,SmId,DateCreated,OofShard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING OrderUid",
		order.TrackNumber,
		order.Entry,
		order.Delivery,
		order.Payment,
		order.Cart,
		order.Locale,
		order.InternalSignature,
		order.CustomerId,
		order.DeliveryService,
		order.Shardkey,
		order.SmId,
		order.DateCreated,
		order.OofShard).Scan(&order.OrderUid)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}
	return order.OrderUid, nil
}
