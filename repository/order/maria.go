package order

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/alek101/GoMikroservisChiNinja/model"
)

func orderIDKey(id uint64) string {
	return fmt.Sprintf("order:%d", id)
}

type MariaRepo struct {
	db *sql.DB
}

func NewMariaRepo(db *sql.DB) *MariaRepo {
	return &MariaRepo{db: db}
}

func (r *MariaRepo) Insert(ctx context.Context, order model.Order) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := orderIDKey(order.OrderID)

	query := `INSERT INTO orders (id, payload) VALUES (?, ?)`
	_, err = r.db.ExecContext(ctx, query, order.OrderID, data)
	if err != nil {
		return fmt.Errorf("failed to insert order %s: %w", key, err)
	}

	return nil
}

func (r *MariaRepo) FindByID(ctx context.Context, id uint64) (model.Order, error) {
	var order model.Order
	if r == nil || r.db == nil {
		return order, fmt.Errorf("database connection not initialized")
	}

	query := `SELECT payload FROM orders WHERE id = ?`
	var data []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return order, fmt.Errorf("order %d not found", id)
		}
		return order, fmt.Errorf("failed to select order %d: %w", id, err)
	}

	if err := json.Unmarshal(data, &order); err != nil {
		return order, fmt.Errorf("failed to decode order payload: %w", err)
	}

	return order, nil
}

// Delete removes an order by its ID.
func (r *MariaRepo) Delete(ctx context.Context, id uint64) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	query := `DELETE FROM orders WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order %d: %w", id, err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("order %d not found", id)
	}
	return nil
}

// Update modifies an existing order record. It overwrites the payload for the
// given OrderID. The provided order object must have a valid OrderID value.
func (r *MariaRepo) Update(ctx context.Context, order model.Order) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	if order.OrderID == 0 {
		return fmt.Errorf("order ID is required for update")
	}

	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := orderIDKey(order.OrderID)

	query := `UPDATE orders SET payload = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, data, order.OrderID)
	if err != nil {
		return fmt.Errorf("failed to update order %s: %w", key, err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("order %s not found", key)
	}

	return nil
}

func (r *MariaRepo) FindAll(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	if r == nil || r.db == nil {
		return orders, fmt.Errorf("database connection not initialized")
	}

	query := `SELECT payload FROM orders`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return orders, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.Order
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return orders, fmt.Errorf("failed to scan order: %w", err)
		}
		if err := json.Unmarshal(data, &order); err != nil {
			return orders, fmt.Errorf("failed to decode order payload: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return orders, fmt.Errorf("error iterating orders: %w", err)
	}

	return orders, nil
}
