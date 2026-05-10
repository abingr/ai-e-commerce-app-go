package repositories

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"ai-e-commerce-app-go/backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) ProductRepository {
	return ProductRepository{db: db}
}

func (r ProductRepository) List(ctx context.Context, filters models.ProductFilters) ([]models.Product, error) {
	query := `
		SELECT id, name, description, brand, category, price_cents, stock_quantity, image_url, is_active, created_at, updated_at
		FROM products
		WHERE is_active = true
	`

	args := []any{}

	if filters.Category != "" {
		args = append(args, filters.Category)
		query += " AND category = $" + argNumber(len(args))
	}

	if filters.Search != "" {
		args = append(args, "%"+strings.ToLower(filters.Search)+"%")
		query += " AND (LOWER(name) LIKE $" + argNumber(len(args)) + " OR LOWER(description) LIKE $" + argNumber(len(args)) + " OR LOWER(brand) LIKE $" + argNumber(len(args)) + ")"
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r ProductRepository) FindByID(ctx context.Context, id string) (models.Product, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, name, description, brand, category, price_cents, stock_quantity, image_url, is_active, created_at, updated_at
		FROM products
		WHERE id = $1 AND is_active = true
	`, id)

	product, err := scanProduct(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Product{}, ErrNotFound
	}

	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (r ProductRepository) Create(ctx context.Context, input models.ProductInput) (models.Product, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO products (name, description, brand, category, price_cents, stock_quantity, image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, name, description, brand, category, price_cents, stock_quantity, image_url, is_active, created_at, updated_at
	`, input.Name, input.Description, input.Brand, input.Category, input.PriceCents, input.StockQuantity, input.ImageURL)

	return scanProduct(row)
}

func (r ProductRepository) Update(ctx context.Context, id string, input models.ProductInput) (models.Product, error) {
	row := r.db.QueryRow(ctx, `
		UPDATE products
		SET name = $2,
			description = $3,
			brand = $4,
			category = $5,
			price_cents = $6,
			stock_quantity = $7,
			image_url = $8,
			updated_at = NOW()
		WHERE id = $1 AND is_active = true
		RETURNING id, name, description, brand, category, price_cents, stock_quantity, image_url, is_active, created_at, updated_at
	`, id, input.Name, input.Description, input.Brand, input.Category, input.PriceCents, input.StockQuantity, input.ImageURL)

	product, err := scanProduct(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.Product{}, ErrNotFound
	}

	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func (r ProductRepository) Delete(ctx context.Context, id string) error {
	commandTag, err := r.db.Exec(ctx, `
		UPDATE products
		SET is_active = false,
			updated_at = NOW()
		WHERE id = $1 AND is_active = true
	`, id)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

type productScanner interface {
	Scan(dest ...any) error
}

func scanProduct(scanner productScanner) (models.Product, error) {
	var product models.Product

	err := scanner.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Brand,
		&product.Category,
		&product.PriceCents,
		&product.StockQuantity,
		&product.ImageURL,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

func argNumber(n int) string {
	return strconv.Itoa(n)
}
