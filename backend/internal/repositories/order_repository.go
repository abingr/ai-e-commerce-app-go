package repositories

import (
	"context"

	"ai-e-commerce-app-go/backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

type orderQuerier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type checkoutItem struct {
	ProductID      string
	ProductName    string
	UnitPriceCents int
	Quantity       int
	LineTotalCents int
}

func NewOrderRepository(db *pgxpool.Pool) OrderRepository {
	return OrderRepository{db: db}
}

func (r OrderRepository) CreateFromCart(ctx context.Context, userID string) (models.Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return models.Order{}, err
	}
	defer tx.Rollback(ctx)

	items, err := r.getCheckoutItems(ctx, tx, userID)
	if err != nil {
		return models.Order{}, err
	}

	if len(items) == 0 {
		return models.Order{}, ErrEmptyCart
	}

	totalCents := 0
	for _, item := range items {
		totalCents += item.LineTotalCents
	}

	var order models.Order
	if err := tx.QueryRow(ctx, `
		INSERT INTO orders (user_id, status, total_cents)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, status, total_cents, created_at, updated_at
	`, userID, models.OrderStatusConfirmed, totalCents).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalCents,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		return models.Order{}, err
	}

	order.Items = make([]models.OrderItem, 0, len(items))
	for _, item := range items {
		var orderItem models.OrderItem
		if err := tx.QueryRow(ctx, `
			INSERT INTO order_items (
				order_id,
				product_id,
				product_name,
				unit_price_cents,
				quantity,
				line_total_cents
			)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, order_id, product_id, product_name, unit_price_cents, quantity, line_total_cents, created_at
		`, order.ID, item.ProductID, item.ProductName, item.UnitPriceCents, item.Quantity, item.LineTotalCents).Scan(
			&orderItem.ID,
			&orderItem.OrderID,
			&orderItem.ProductID,
			&orderItem.ProductName,
			&orderItem.UnitPriceCents,
			&orderItem.Quantity,
			&orderItem.LineTotalCents,
			&orderItem.CreatedAt,
		); err != nil {
			return models.Order{}, err
		}

		order.Items = append(order.Items, orderItem)
	}

	if _, err := tx.Exec(ctx, `
		DELETE FROM cart_items
		WHERE user_id = $1
	`, userID); err != nil {
		return models.Order{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (r OrderRepository) ListByUser(ctx context.Context, userID string) ([]models.Order, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, status, total_cents, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []models.Order{}
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalCents,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, err
		}

		items, err := r.loadOrderItems(ctx, r.db, order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r OrderRepository) FindByIDForUser(ctx context.Context, userID string, orderID string) (models.Order, error) {
	var order models.Order
	if err := r.db.QueryRow(ctx, `
		SELECT id, user_id, status, total_cents, created_at, updated_at
		FROM orders
		WHERE id = $1 AND user_id = $2
	`, orderID, userID).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalCents,
		&order.CreatedAt,
		&order.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			return models.Order{}, ErrNotFound
		}
		return models.Order{}, err
	}

	items, err := r.loadOrderItems(ctx, r.db, order.ID)
	if err != nil {
		return models.Order{}, err
	}
	order.Items = items

	return order, nil
}

func (r OrderRepository) getCheckoutItems(ctx context.Context, q orderQuerier, userID string) ([]checkoutItem, error) {
	rows, err := q.Query(ctx, `
		SELECT
			p.id,
			p.name,
			p.price_cents,
			ci.quantity,
			ci.quantity * p.price_cents AS line_total_cents
		FROM cart_items ci
		JOIN products p ON p.id = ci.product_id
		WHERE ci.user_id = $1 AND p.is_active = true
		ORDER BY ci.created_at ASC
		FOR UPDATE OF ci
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []checkoutItem{}
	for rows.Next() {
		var item checkoutItem
		if err := rows.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.UnitPriceCents,
			&item.Quantity,
			&item.LineTotalCents,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r OrderRepository) loadOrderItems(ctx context.Context, q orderQuerier, orderID string) ([]models.OrderItem, error) {
	rows, err := q.Query(ctx, `
		SELECT id, order_id, product_id, product_name, unit_price_cents, quantity, line_total_cents, created_at
		FROM order_items
		WHERE order_id = $1
		ORDER BY created_at ASC
	`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.OrderItem{}
	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.ProductName,
			&item.UnitPriceCents,
			&item.Quantity,
			&item.LineTotalCents,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
