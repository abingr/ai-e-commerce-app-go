package repositories

import (
	"context"
	"errors"
	"strings"

	"ai-e-commerce-app-go/backend/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, password_hash, role, created_at, updated_at
	`, user.Name, strings.ToLower(user.Email), user.PasswordHash, user.Role)

	createdUser, err := scanUser(row)
	if isUniqueViolation(err) {
		return models.User{}, ErrConflict
	}

	if err != nil {
		return models.User{}, err
	}

	return createdUser, nil
}

func (r UserRepository) FindByEmail(ctx context.Context, email string) (models.User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, name, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`, strings.ToLower(email))

	user, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrNotFound
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r UserRepository) FindByID(ctx context.Context, id string) (models.User, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, name, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id)

	user, err := scanUser(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrNotFound
	}

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

type userScanner interface {
	Scan(dest ...any) error
}

func scanUser(scanner userScanner) (models.User, error) {
	var user models.User

	err := scanner.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
