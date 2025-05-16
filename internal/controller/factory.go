package controller

import (
	"fmt"
	"strings"

	"github.com/go-faker/faker/v4"
	"github.com/vvwind/2025-MAI-Backend-V-Vetrov/internal/model"
)

type UserFactory struct {
	ID       int64
	UserName string
	Email    string
	Role     string
}

func (f UserFactory) Build() *model.User {
	user := &model.User{
		ID:       f.ID,
		UserName: f.UserName,
		Email:    f.Email,
		Role:     f.Role,
	}

	if user.UserName == "" {
		user.UserName = strings.Replace(faker.Username(), " ", "", -1)
	}
	if user.Email == "" {
		user.Email = faker.Email()
	}
	if user.Role == "" {
		user.Role = "customer"
	}
	if user.ID == 0 {
		user.ID = faker.RandomUnixTime() % 10000
	}

	return user
}

type UserRegisterFactory struct {
	UserName string
	Email    string
	Password string
	Role     string
}

func (f UserRegisterFactory) Build() *model.UserRegister {
	user := &model.UserRegister{
		UserName: f.UserName,
		Email:    f.Email,
		Password: f.Password,
		Role:     f.Role,
	}

	if user.UserName == "" {
		user.UserName = strings.Replace(faker.Username(), " ", "", -1)
	}
	if user.Email == "" {
		user.Email = faker.Email()
	}
	if user.Password == "" {
		user.Password = "password" + fmt.Sprint(faker.RandomUnixTime()%1000)
	}
	if user.Role == "" {
		user.Role = "customer"
	}

	return user
}

type UserLoginFactory struct {
	Email    string
	Password string
}

func (f UserLoginFactory) Build() *model.UserLogin {
	login := &model.UserLogin{
		Email:    f.Email,
		Password: f.Password,
	}

	if login.Email == "" {
		login.Email = faker.Email()
	}
	if login.Password == "" {
		login.Password = "password" + fmt.Sprint(faker.RandomUnixTime()%1000)
	}

	return login
}

type ProductFactory struct {
	ID                 int64
	Title              string
	SellerName         string
	SellerID           int64
	ProductDescription string
	ProductImage       string
	Price              int64
	Amount             int
}

func (f ProductFactory) Build() *model.Product {
	product := &model.Product{
		ID:                 f.ID,
		Title:              f.Title,
		SellerName:         f.SellerName,
		SellerID:           f.SellerID,
		ProductDescription: f.ProductDescription,
		ProductImage:       f.ProductImage,
		Price:              f.Price,
		Amount:             f.Amount,
	}

	if product.Title == "" {
		product.Title = faker.Word()
	}
	if product.SellerName == "" {
		product.SellerName = strings.Replace(faker.Name(), " ", "", -1)
	}
	if product.ProductDescription == "" {
		product.ProductDescription = faker.Sentence()
	}
	if product.ProductImage == "" {
		product.ProductImage = fmt.Sprintf("%s.jpg", faker.Word())
	}
	if product.Price == 0 {
		product.Price = faker.RandomUnixTime() % 1000
	}
	if product.Amount == 0 {
		product.Amount = int(faker.RandomUnixTime() % 100)
	}
	if product.SellerID == 0 {
		product.SellerID = faker.RandomUnixTime() % 10000
	}
	if product.ID == 0 {
		product.ID = faker.RandomUnixTime() % 10000
	}

	return product
}

type CreateProductRequestFactory struct {
	Title              string
	ProductDescription string
	ProductImage       string
	Price              int64
	Amount             int
}

func (f CreateProductRequestFactory) Build() *model.CreateProductRequest {
	req := &model.CreateProductRequest{
		Title:              f.Title,
		ProductDescription: f.ProductDescription,
		ProductImage:       f.ProductImage,
		Price:              f.Price,
		Amount:             f.Amount,
	}

	if req.Title == "" {
		req.Title = faker.Word()
	}
	if req.ProductDescription == "" {
		req.ProductDescription = faker.Sentence()
	}
	if req.ProductImage == "" {
		req.ProductImage = fmt.Sprintf("%s.jpg", faker.Word())
	}
	if req.Price == 0 {
		req.Price = faker.RandomUnixTime() % 1000
	}
	if req.Amount == 0 {
		req.Amount = int(faker.RandomUnixTime() % 100)
	}

	return req
}

type UpdateProductRequestFactory struct {
	Title              string
	ProductDescription string
	ProductImage       string
	Price              int64
	Amount             int
}

func (f UpdateProductRequestFactory) Build() *model.UpdateProductRequest {
	req := &model.UpdateProductRequest{
		Title:              f.Title,
		ProductDescription: f.ProductDescription,
		ProductImage:       f.ProductImage,
		Price:              f.Price,
		Amount:             f.Amount,
	}

	if req.Title == "" {
		req.Title = "Updated " + faker.Word()
	}
	if req.ProductDescription == "" {
		req.ProductDescription = "Updated " + faker.Sentence()
	}
	if req.ProductImage == "" {
		req.ProductImage = "updated_" + faker.Word() + ".jpg"
	}
	if req.Price == 0 {
		req.Price = faker.RandomUnixTime() % 1000
	}
	if req.Amount == 0 {
		req.Amount = int(faker.RandomUnixTime() % 100)
	}

	return req
}
