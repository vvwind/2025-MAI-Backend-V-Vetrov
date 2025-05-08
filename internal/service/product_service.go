package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/vvwind/2025-MAI-Backend-V-Vetrov/internal/model"
	"github.com/vvwind/2025-MAI-Backend-V-Vetrov/internal/repository"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock
type ProductService interface {
	GetAllProducts(ctx context.Context) ([]model.Product, error)
	GetProductByID(ctx context.Context, id int64) (*model.Product, error)
	CreateProduct(ctx context.Context, ProductReq model.CreateProductRequest, seller model.User) (int64, error)
	UpdateProduct(ctx context.Context, productReq model.UpdateProductRequest, productID, userID int64) (int64, error)
	DeleteProduct(ctx context.Context, id int64) error
	AddToCart(ctx context.Context, productID, userID int64) error
	BuyProduct(ctx context.Context, productID, userID int64) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAllProducts(ctx context.Context) ([]model.Product, error) {
	return s.repo.GetAllProducts(ctx)
}

func (s *productService) GetProductByID(ctx context.Context, id int64) (*model.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *productService) CreateProduct(ctx context.Context, ProductReq model.CreateProductRequest, seller model.User) (int64, error) {

	newProduct := ConvertRequestToProduct(ProductReq)
	newProduct.SellerID = seller.ID
	newProduct.SellerName = seller.UserName
	return s.repo.CreateProduct(ctx, newProduct)
}

func (s *productService) UpdateProduct(ctx context.Context, productReq model.UpdateProductRequest,
	productID, userID int64) (int64, error) {

	existingProduct, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return -1, fmt.Errorf("error getting product data: %w", err)
	}

	if existingProduct == nil {
		return -1, fmt.Errorf("product does not exist")
	}

	if existingProduct.SellerID != userID {
		return -1, errors.New("seller ID and product ID does not match")
	}

	query := "UPDATE products SET "
	params := []interface{}{}
	paramCount := 1

	updates := []string{}

	if productReq.Title != "" {
		updates = append(updates, fmt.Sprintf("title = $%d", paramCount))
		params = append(params, productReq.Title)
		paramCount++
	}

	if productReq.ProductDescription != "" {
		updates = append(updates, fmt.Sprintf("product_description = $%d", paramCount))
		params = append(params, productReq.ProductDescription)
		paramCount++
	}

	if productReq.ProductImage != "" {
		updates = append(updates, fmt.Sprintf("product_image = $%d", paramCount))
		params = append(params, productReq.ProductImage)
		paramCount++
	}

	if productReq.Price != 0 {
		updates = append(updates, fmt.Sprintf("price = $%d", paramCount))
		params = append(params, productReq.Price)
		paramCount++
	}

	if productReq.Amount != 0 {
		updates = append(updates, fmt.Sprintf("amount = $%d", paramCount))
		params = append(params, productReq.Amount)
		paramCount++
	}

	if len(updates) == 0 {
		return -1, errors.New("nothing to update")
	}

	query += strings.Join(updates, ", ")
	query += ", updated_at = NOW() WHERE id = $" + strconv.Itoa(paramCount)
	params = append(params, productID)

	query += " RETURNING id;"

	return s.repo.UpdateProduct(ctx, query, params)
}

func (s *productService) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.DeleteProduct(ctx, id)
}

func (s *productService) AddToCart(ctx context.Context, productID, userID int64) error {
	return s.repo.AddToCart(ctx, productID, userID)
}

func (s *productService) BuyProduct(ctx context.Context, productID, userID int64) error {
	err := s.repo.CheckCart(ctx, productID, userID)

	if err != nil {
		return err
	}

	return s.repo.BuyProduct(ctx, productID, userID)
}
