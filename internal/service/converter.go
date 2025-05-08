package service

import "github.com/vvwind/2025-MAI-Backend-V-Vetrov/internal/model"

func FromRequestToModel(usr model.UserRegister) model.User {
	return model.User{UserName: usr.UserName, Email: usr.Email, Role: usr.Role}
}

func ConvertRequestToProduct(req model.CreateProductRequest) model.Product {
	return model.Product{
		Title:              req.Title,
		ProductDescription: req.ProductDescription,
		ProductImage:       req.ProductImage,
		Price:              req.Price,
	}
}
