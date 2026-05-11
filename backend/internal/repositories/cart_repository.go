package repositories

import (
	"context"

	"ai-e-commerce-app-go/backend/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) CartRepository {
	return CartRepository{db: db}
}

func (r CartRepository) Get(ctx context.Context, userID string) (models.Cart, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			p.id,
			p.name,
			p.brand,
			p.category,
			p.image_url,
			ci.quantity,
			p.price_cents,
			ci.quantity * p.price_cents AS line_total_cents
		FROM cart_items ci
		JOIN products p ON p.id = ci.product_id
		WHERE ci.user_id = $1 AND p.is_active = true
		ORDER BY ci.created_at DESC
	`, userID)
	if err != nil {
		return models.Cart{}, err
	}
	defer rows.Close()

	cart := models.Cart{
		Items: []models.CartItem{},
	}

	for rows.Next() {
		var item models.CartItem
		if err := rows.Scan(
			&item.ProductID,
			&item.Name,
			&item.Brand,
			&item.Category,
			&item.ImageURL,
			&item.Quantity,
			&item.UnitPriceCents,
			&item.LineTotalCents,
		); err != nil {
			return models.Cart{}, err
		}

		cart.Items = append(cart.Items, item)
		cart.TotalCents += item.LineTotalCents
	}

	if err := rows.Err(); err != nil {
		return models.Cart{}, err
	}

	return cart, nil
}

func (r CartRepository) AddItem(ctx context.Context, userID string, input models.AddCartItemInput) error {
	commandTag, err := r.db.Exec(ctx, `
		INSERT INTO cart_items (user_id, product_id, quantity)
		SELECT $1, p.id, $3
		FROM products p
		WHERE p.id = $2 AND p.is_active = true
		ON CONFLICT (user_id, product_id)
		DO UPDATE SET
			quantity = cart_items.quantity + EXCLUDED.quantity,
			updated_at = NOW()
	`, userID, input.ProductID, input.Quantity)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r CartRepository) UpdateItem(ctx context.Context, userID string, productID string, quantity int) error {
	commandTag, err := r.db.Exec(ctx, `
		UPDATE cart_items ci
		SET quantity = $3,
			updated_at = NOW()
		FROM products p
		WHERE ci.product_id = p.id
			AND ci.user_id = $1
			AND ci.product_id = $2
			AND p.is_active = true
	`, userID, productID, quantity)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r CartRepository) RemoveItem(ctx context.Context, userID string, productID string) error {
	commandTag, err := r.db.Exec(ctx, `
		DELETE FROM cart_items
		WHERE user_id = $1 AND product_id = $2
	`, userID, productID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r CartRepository) Clear(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM cart_items
		WHERE user_id = $1
	`, userID)
	return err
}
