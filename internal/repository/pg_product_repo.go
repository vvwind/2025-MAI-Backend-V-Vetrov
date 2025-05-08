package repository

import (
	"context"
	"fmt"
	"rest-api/internal/model"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	cartKey = "cart"
)

type ProductRepository interface {
	GetAllProducts(ctx context.Context) ([]model.Product, error)
	GetProductByID(ctx context.Context, id int64) (*model.Product, error)
	CreateProduct(ctx context.Context, product model.Product) (int64, error)
	UpdateProduct(ctx context.Context, query string, params []interface{}) (int64, error)
	DeleteProduct(ctx context.Context, id int64) error
	CheckAccess(ctx context.Context, productID int64) (int64, error)
	AddToCart(ctx context.Context, productID int64, userID int64) error
	CheckCart(ctx context.Context, productID int64, userID int64) error
	BuyProduct(ctx context.Context, productID int64, userID int64) error
}

type postgresProductRepository struct {
	pool *pgxpool.Pool
	rc   *redis.Client
}

func NewPostgresProductRepository(pool *pgxpool.Pool, rc *redis.Client) ProductRepository {
	return &postgresProductRepository{pool: pool}
}

func (r *postgresProductRepository) GetAllProducts(ctx context.Context) ([]model.Product, error) {
	query := `SELECT id,
	title, 
	seller_name, 
	seller_id, 
	product_description, 
	product_image, 
	price, 
	amount
	FROM products;`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.SellerName,
			&p.SellerID,
			&p.ProductDescription,
			&p.ProductImage,
			&p.Price,
			&p.Amount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return products, nil
}

func (r *postgresProductRepository) GetProductByID(ctx context.Context, id int64) (*model.Product, error) {
	query := `SELECT id,
	title, 
	seller_name, 
	seller_id, 
	product_description, 
	product_image, 
	price, 
	amount
	FROM products
	WHERE id = $1;`
	row := r.pool.QueryRow(ctx, query, id)

	var p model.Product
	err := row.Scan(
		&p.ID,
		&p.Title,
		&p.SellerName,
		&p.SellerID,
		&p.ProductDescription,
		&p.ProductImage,
		&p.Price,
		&p.Amount,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &p, nil
}

func (r *postgresProductRepository) CreateProduct(ctx context.Context, product model.Product) (int64, error) {
	query := `
		INSERT INTO products 
		(title, seller_name, seller_id, product_image, 
		product_description, price, amount, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW()) 
		RETURNING id;
	`
	row := r.pool.QueryRow(
		ctx,
		query,
		product.Title,
		product.SellerName,
		product.SellerID,
		product.ProductImage,
		product.ProductDescription,
		product.Price,
		product.Amount,
	)

	var createdID int64
	err := row.Scan(
		&createdID,
	)

	if err != nil {
		return -1, fmt.Errorf("failed to create product: %w", err)
	}

	return createdID, nil
}

func (r *postgresProductRepository) UpdateProduct(ctx context.Context,
	query string, params []interface{}) (int64, error) {

	row := r.pool.QueryRow(
		ctx,
		query,
		params...,
	)

	var updatedID int64
	err := row.Scan(
		&updatedID,
	)

	if err != nil {
		return -1, fmt.Errorf("failed to update product: %w", err)
	}

	return updatedID, nil
}

func (r *postgresProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1;`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

func (r *postgresProductRepository) AddToCart(ctx context.Context, productID int64, userID int64) error {
	err := r.rc.Set(ctx, fmt.Sprintf("%s_%d_%d", cartKey, userID, productID), "1", time.Hour)
	if err.Err() != nil {
		return fmt.Errorf("error setting redis key: %w", err.Err())
	}

	return nil
}

func (r *postgresProductRepository) CheckCart(ctx context.Context, productID int64, userID int64) error {
	err := r.rc.Get(ctx, fmt.Sprintf("%s_%d_%d", cartKey, userID, productID))
	if err.Err() != nil {
		return fmt.Errorf("error getting redis key: %w", err.Err())
	}

	return nil
}

func (r *postgresProductRepository) CheckAccess(ctx context.Context, productID int64) (int64, error) {
	query := `SELECT seller_id FROM products WHERE id = $1;`
	row := r.pool.QueryRow(ctx, query, productID)
	var sellerID int64
	err := row.Scan(&sellerID)

	if err != nil {
		return -1, fmt.Errorf("error checking access: %w", err)
	}

	return sellerID, nil
}

func (r *postgresProductRepository) BuyProduct(ctx context.Context, productID int64, userID int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var currentAmount int
	err = tx.QueryRow(ctx,
		"SELECT amount FROM products WHERE id = $1 FOR UPDATE",
		productID).Scan(&currentAmount)

	if err != nil {
		return fmt.Errorf("failed to query product amount: %w", err)
	}

	if currentAmount <= 0 {
		return fmt.Errorf("product out of stock")
	}

	_, err = tx.Exec(ctx,
		"UPDATE products SET amount = amount - 1 WHERE id = $1",
		productID)
	if err != nil {
		return fmt.Errorf("failed to update product amount: %w", err)
	}

	rowRedis := r.rc.Del(ctx, fmt.Sprintf("%s_%d_%d", cartKey, userID, productID))

	if rowRedis.Err() != nil {
		return fmt.Errorf("failed to delete redis key: %w", rowRedis.Err())
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
