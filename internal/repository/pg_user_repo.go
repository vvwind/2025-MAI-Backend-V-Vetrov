package repository

import (
	"context"
	"fmt"

	"github.com/vvwind/2025-MAI-Backend-V-Vetrov/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, usr model.User, passwordHash string) (int64, error)
	GetHashedPassword(ctx context.Context, email string) (string, error)
}

type postgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) UserRepository {
	return &postgresUserRepository{pool: pool}
}

func (r *postgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, user_name, email, user_role FROM users WHERE email = $1;`
	row := r.pool.QueryRow(ctx, query, email)
	var usr model.User
	err := row.Scan(&usr.ID, &usr.UserName, &usr.Email, &usr.Role)
	if err != nil {
		return nil, fmt.Errorf("error performing get user query: %w", err)
	}

	return &usr, nil
}

func (r *postgresUserRepository) GetHashedPassword(ctx context.Context, email string) (string, error) {
	query := `SELECT password_hash FROM users_creds WHERE email = $1;`
	row := r.pool.QueryRow(ctx, query, email)
	var hashedPass string
	err := row.Scan(&hashedPass)
	if err != nil {
		return "", fmt.Errorf("error scanning hashed password: %w", err)
	}

	return hashedPass, nil
}

func (r *postgresUserRepository) CreateUser(ctx context.Context, usr model.User, passwordHash string) (int64, error) {
	// Begin transaction
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. First insert into users table (parent table)
	var userID int64
	userQuery := `INSERT INTO users (user_name, email, user_role, created_at, updated_at)
                 VALUES ($1, $2, $3, NOW(), NOW())
                 RETURNING id`

	err = tx.QueryRow(ctx, userQuery, usr.UserName, usr.Email, usr.Role).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	// 2. Then insert into users_creds table (child table)
	credsQuery := `INSERT INTO users_creds (user_id, email, password_hash, created_at, updated_at)
                  VALUES ($1, $2, $3, NOW(), NOW())`

	_, err = tx.Exec(ctx, credsQuery, userID, usr.Email, passwordHash)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user credentials: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, nil
}
