package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/popeskul/payment-gateway/internal/core/domain/user"
	"github.com/popeskul/payment-gateway/internal/core/ports"
)

type UserRepository struct {
	db            *Database
	uuidGenerator ports.UUIDGenerator
}

func NewUserRepository(db *Database, uuidGenerator ports.UUIDGenerator) ports.UserRepository {
	return &UserRepository{
		db:            db,
		uuidGenerator: uuidGenerator,
	}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	if u.ID == "" {
		u.ID = r.uuidGenerator.Generate()
	}

	query := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		u.ID, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var u user.User
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var u user.User
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	query := `
		UPDATE users
		SET email = $2, password_hash = $3, first_name = $4, last_name = $5, updated_at = $6
		WHERE id = $1
	`
	_, err := r.db.Pool.Exec(ctx, query,
		u.ID, u.Email, u.PasswordHash, u.FirstName, u.LastName, u.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}
