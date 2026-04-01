package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/online-shop/internal/models"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create inserts an order and its items inside a transaction.
func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	orderQuery := `
		INSERT INTO orders (user_id, status, total_price, address)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	err = tx.QueryRow(ctx, orderQuery,
		order.UserID, order.Status, order.TotalPrice, order.Address,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	itemQuery := `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	for i := range order.Items {
		order.Items[i].OrderID = order.ID
		err = tx.QueryRow(ctx, itemQuery,
			order.Items[i].OrderID, order.Items[i].ProductID,
			order.Items[i].Quantity, order.Items[i].Price,
		).Scan(&order.Items[i].ID)
		if err != nil {
			return fmt.Errorf("insert order item: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*models.Order, error) {
	var o models.Order
	orderQuery := `SELECT id, user_id, status, total_price, address, created_at, updated_at
	               FROM orders WHERE id = $1`

	err := r.db.QueryRow(ctx, orderQuery, id).Scan(
		&o.ID, &o.UserID, &o.Status, &o.TotalPrice,
		&o.Address, &o.CreatedAt, &o.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	itemQuery := `SELECT id, order_id, product_id, quantity, price
	              FROM order_items WHERE order_id = $1`
	rows, err := r.db.Query(ctx, itemQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, item)
	}

	return &o, nil
}

func (r *OrderRepository) ListByUserID(ctx context.Context, userID int64) ([]models.Order, error) {
	query := `SELECT id, user_id, status, total_price, address, created_at, updated_at
	           FROM orders WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.TotalPrice, &o.Address, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	return orders, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, id int64, status models.OrderStatus) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	ct, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
