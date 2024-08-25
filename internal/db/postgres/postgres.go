package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"time"
	"wb/internal/models"
)

type RepoPostgres struct {
	connection *pgxpool.Pool
}

func Connect(ctx context.Context, url string) (*RepoPostgres, error) {
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при парсинге конфига базы данных: %w", err)
	}

	dbConn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %w", err)
	}

	conn, err := dbConn.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("ошибка получения соединения: %w", err)
	} else {
		log.Println("Соединение с базой данных установлено")
	}
	defer conn.Release()

	time.Sleep(1 * time.Second)

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
		return 0, fmt.Errorf("%w", err)
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
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee)
	if err != nil {
		return 0, fmt.Errorf("failed to create payments table: %w", err)
	}

	for _, item := range order.Items {
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

func (repo *RepoPostgres) GetOrder(ctx context.Context, orderUID string) bool {

	var count int
	err := repo.connection.QueryRow(ctx, "SELECT COUNT(*) FROM orders WHERE order_uid = $1", orderUID).Scan(&count)
	if err != nil {
		log.Errorf("Ошибка при соединении с бд: %s", err)
		return false
	}

	if count == 0 {
		return false
	} else {
		return true
	}
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
			&order.ID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerId,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SmId,
			&order.DateCreated,
			&order.OofShard)
		if err != nil {
			return nil, err
		}

		order.Delivery, err = GetDeliveryByOrderUID(ctx, repo.connection, order.OrderUid)
		if err != nil {
			log.Fatalf("Ошибка при формировании структуры delivery: %v", err)
		}
		order.Payment, err = GetPaymentByOrderUID(ctx, repo.connection, order.OrderUid)
		if err != nil {
			log.Fatalf("Ошибка при формировании структуры payment: %v", err)
		}
		order.Items, err = GetItemsByOrderUID(ctx, repo.connection, order.OrderUid)
		if err != nil {
			log.Fatalf("Ошибка при формировании структуры items: %v", err)
		}

		orders = append(orders, order)
	}
	defer rows.Close()
	return orders, err
}

func GetDeliveryByOrderUID(ctx context.Context, pool *pgxpool.Pool, orderUID string) (models.Delivery, error) {
	var delivery models.Delivery

	err := pool.QueryRow(ctx, "SELECT * FROM deliveries WHERE order_uid = $1", orderUID).Scan(
		&delivery.OrderUid,
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email)
	if err != nil {
		return delivery, err
	}
	return delivery, nil
}

func GetPaymentByOrderUID(ctx context.Context, pool *pgxpool.Pool, orderUID string) (models.Payment, error) {
	var payment models.Payment
	err := pool.QueryRow(ctx, "SELECT * FROM payments WHERE order_uid = $1", orderUID).Scan(
		&payment.OrderUid,
		&payment.Transaction,
		&payment.RequestId,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)
	if err != nil {
		return payment, err
	}
	return payment, nil
}

func GetItemsByOrderUID(ctx context.Context, pool *pgxpool.Pool, orderUID string) ([]models.Items, error) {
	rows, err := pool.Query(ctx, "SELECT * FROM items WHERE order_uid = $1", orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Items, 0)
	for rows.Next() {
		var item models.Items
		err = rows.Scan(
			&item.OrderUid,
			&item.ChrtId,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmId,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
