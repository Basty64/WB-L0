package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
	"wb/internal/models"
)

type RepoPostgres struct {
	connection *pgxpool.Pool
}

func (repo *RepoPostgres) Deadline() (deadline time.Time, ok bool) {
	//TODO implement me
	panic("implement me")
}

func (repo *RepoPostgres) Done() <-chan struct{} {
	//TODO implement me
	panic("implement me")
}

func (repo *RepoPostgres) Err() error {
	//TODO implement me
	panic("implement me")
}

func (repo *RepoPostgres) Value(key any) any {
	//TODO implement me
	panic("implement me")
}

func Connect(ctx context.Context, url string) (*RepoPostgres, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatalf("failed to parse the config: %w", err)
	}

	dbConn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("failed to connect to database: %w", err)
	}

	return &RepoPostgres{connection: dbConn}, nil
}

func (repo *RepoPostgres) CloseConnection() {
	repo.connection.Close()
}

func (repo *RepoPostgres) InsertOrder(ctx context.Context, order models.Order) (int, error) {

	err := repo.connection.QueryRow(ctx, "INSERT INTO orders ("+
		"TrackNumber, Entry, Delivery, Payment, Cart, Locale, InternalSignature, CustomerId, DeliveryService, Shardkey, SmId, DateCreated, OofShard) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING OrderUid",
		order.TrackNumber,
		order.Entry,
		order.Delivery,
		order.Payment,
		order.Item,
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

func (repo *RepoPostgres) GetOrder(ctx context.Context, orderUid int) (models.Order, error) {
	var order models.Order

	row, err := repo.connection.Query(ctx, "SELECT * FROM orders WHERE id = $1", orderUid)
	if err != nil {
		return models.Order{}, err
	}

	for row.Next() {
		err := row.Scan(
			&order.OrderUid,
			&order.TrackNumber,
			&order.Entry,
			&order.Delivery,
			&order.Payment,
			&order.Item,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerId,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SmId,
			&order.DateCreated,
			&order.OofShard)
		if err != nil {
			return models.Order{}, err
		}
	}
	defer row.Close()
	return order, err
}

func (repo *RepoPostgres) GetAllOrders(ctx context.Context) ([]models.Order, error) {

	var orders []models.Order

	rows, err := repo.connection.Query(ctx, "SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.OrderUid,
			&order.TrackNumber,
			&order.Entry,
			&order.Delivery,
			&order.Payment,
			&order.Item,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerId,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SmId)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	defer rows.Close()
	return orders, err
}
