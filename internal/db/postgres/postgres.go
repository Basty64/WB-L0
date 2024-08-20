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

func (repo *RepoPostgres) InsertOrder(ctx context.Context, order *models.Order) (int, error) {

	err := repo.connection.QueryRow(ctx, "INSERT INTO orders (order_uid, track_number, entry, Locale,internal_signature, customer_id, delivery_service, Shardkey, sm_id, date_created, oofshard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
		order.OrderUid,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerId,
		order.DeliveryService,
		order.Shardkey,
		order.SmId,
		order.DateCreated,
		order.OofShard).Scan(&order.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	_, err = repo.connection.Exec(ctx, "INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		order.OrderUid,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email)
	if err != nil {
		return 0, fmt.Errorf("failed to create deliveries table: %w", err)
	}

	_, err = repo.connection.Exec(ctx, "INSERT INTO payments (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		order.OrderUid,
		order.Payment.Transaction,
		order.Payment.RequestId,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		time.Unix(int64(order.Payment.PaymentDt), 0),
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee)
	if err != nil {
		return 0, fmt.Errorf("failed to create payments table: %w", err)
	}

	for _, item := range order.Item {
		_, err = repo.connection.Exec(ctx, "INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
			order.OrderUid,
			item.ChrtId,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmId,
			item.Brand,
			item.Status)
		if err != nil {
			return 0, fmt.Errorf("failed to create items table: %w", err)
		}
	}
	return order.ID, nil
}

func (repo *RepoPostgres) GetOrder(ctx context.Context, ID int) (models.Order, error) {
	var order models.Order

	row, err := repo.connection.Query(ctx, "SELECT * FROM orders WHERE id = $1", ID)
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
