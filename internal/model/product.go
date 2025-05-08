package model

type User struct {
	ID       int64  `json:"id"`
	UserName string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type UserRegister struct {
	UserName string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Product struct {
	ID                 int64  `json:"id"`
	Title              string `json:"title"`
	SellerName         string `json:"seller_name"`
	SellerID           int64  `json:"seller_id"`
	ProductDescription string `json:"product_description"`
	ProductImage       string `json:"product_image"`
	Price              int64  `json:"price"`
	Amount             int    `json:"amount"`
}

type CreateProductRequest struct {
	Title              string `json:"title"`
	ProductDescription string `json:"product_description"`
	ProductImage       string `json:"product_image"`
	Price              int64  `json:"price"`
	Amount             int    `json:"amount"`
}

type UpdateProductRequest struct {
	Title              string `json:"title"`
	ProductDescription string `json:"product_description"`
	ProductImage       string `json:"product_image"`
	Price              int64  `json:"price"`
	Amount             int    `json:"amount"`
}

type A struct{}
