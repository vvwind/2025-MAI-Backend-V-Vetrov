package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/vvwind/2025-MAI-Backend-V-Vetrov/internal/middleware"
	"github.com/vvwind/2025-MAI-Backend-V-Vetrov/internal/model"
	"github.com/vvwind/2025-MAI-Backend-V-Vetrov/internal/service"
	"github.com/vvwind/2025-MAI-Backend-V-Vetrov/pkg/utils"

	"github.com/gorilla/mux"
)

type MarketplaceController struct {
	prSrvc  service.ProductService
	usrSrvc service.UserService
}

func NewMarketplaceController(servicePr service.ProductService, serviceUs service.UserService) *MarketplaceController {
	return &MarketplaceController{
		prSrvc:  servicePr,
		usrSrvc: serviceUs,
	}
}

func (c *MarketplaceController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/user", c.CreateUser).Methods("POST")
	router.HandleFunc("/user/login", c.LoginUser).Methods("POST")

	// Protected routes (auth required)
	protectedRouter := router.PathPrefix("").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware)

	protectedRouter.HandleFunc("/products", c.GetAllProducts).Methods("GET")
	protectedRouter.HandleFunc("/products/{id}", c.GetProductByID).Methods("GET")
	protectedRouter.HandleFunc("/products", c.CreateProduct).Methods("POST")
	protectedRouter.HandleFunc("/products/{id}", c.UpdateProduct).Methods("PUT")
	protectedRouter.HandleFunc("/products/{id}", c.DeleteProduct).Methods("DELETE")

	protectedRouter.HandleFunc("/products/cart/{id}", c.AddToCart).Methods("POST")
	protectedRouter.HandleFunc("/products/buy/{id}", c.BuyProduct).Methods("POST")
}

func (c *MarketplaceController) CreateUser(w http.ResponseWriter, r *http.Request) {

	var err error

	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	ctx, cancel := context.WithTimeout(r.Context(), 50*time.Second)
	defer cancel()

	// Parse request body
	var req model.UserRegister
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate required fields
	if req.UserName == "" || req.Email == "" || req.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Name, email and password are required")
		return
	}

	if req.Role != "customer" && req.Role != "seller" {
		utils.RespondWithError(w, http.StatusBadRequest, "Wrong role")
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	// Call service layer
	userID, err := c.usrSrvc.CreateUser(ctx, req, hashedPassword)
	if err != nil {
		// Handle specific errors if needed
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Return response (omitting password hash)
	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"id":       userID,
		"username": req.UserName,
		"email":    req.Email,
		"role":     req.Role,
	})
}

func (c *MarketplaceController) LoginUser(w http.ResponseWriter, r *http.Request) {

	var err error

	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	ctx, cancel := context.WithTimeout(r.Context(), 50*time.Second)
	defer cancel()

	var loginReq model.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate input
	if loginReq.Email == "" || loginReq.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Call service
	token, err := c.usrSrvc.LoginUser(ctx, loginReq)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"TokenBearer": token,
	})
}

func (c *MarketplaceController) GetAllProducts(w http.ResponseWriter, r *http.Request) {

	var err error

	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	ctx, cancel := context.WithTimeout(r.Context(), 50*time.Second)
	defer cancel()

	products, err := c.prSrvc.GetAllProducts(ctx)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, products)
}

func (c *MarketplaceController) GetProductByID(w http.ResponseWriter, r *http.Request) {

	var err error

	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	ctx, cancel := context.WithTimeout(r.Context(), 50*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid id")
		return
	}

	Product, err := c.prSrvc.GetProductByID(ctx, intId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if Product == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, Product)
}

func (c *MarketplaceController) CreateProduct(w http.ResponseWriter, r *http.Request) {

	var err error

	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	ctx, cancel := context.WithTimeout(r.Context(), 50*time.Second)
	defer cancel()

	var prReq model.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&prReq); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if prReq.Title == "" || prReq.ProductImage == "" || prReq.Price <= 0 || prReq.Amount <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Wrong product data")
		return
	}

	claims, ok := utils.GetUserClaimsFromContext(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user claims")
		return
	}

	userEmail, ok := claims["email"].(string)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "User email not found in token")
		return
	}

	curUser, err := c.usrSrvc.GetUserByEmail(ctx, userEmail)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "User not found by email")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&prReq); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	product, err := c.prSrvc.CreateProduct(ctx, prReq, *curUser)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, product)
}

func (c *MarketplaceController) AddToCart(w http.ResponseWriter, r *http.Request) {

	var err error

	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	claims, ok := utils.GetUserClaimsFromContext(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user claims")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product id")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 50*time.Second)
	defer cancel()

	userEmail, ok := claims["email"].(string)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "User email not found in token")
		return
	}

	curUser, err := c.usrSrvc.GetUserByEmail(ctx, userEmail)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "User not found by email")
		return
	}

	c.prSrvc.AddToCart(ctx, intId, curUser.ID)

}

func (c *MarketplaceController) BuyProduct(w http.ResponseWriter, r *http.Request) {

	var err error

	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	claims, ok := utils.GetUserClaimsFromContext(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user claims")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product id")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 50*time.Second)
	defer cancel()

	userEmail, ok := claims["email"].(string)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "User email not found in token")
		return
	}

	curUser, err := c.usrSrvc.GetUserByEmail(ctx, userEmail)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "User not found by email")
		return
	}

	c.prSrvc.BuyProduct(ctx, intId, curUser.ID)

}

func (c *MarketplaceController) UpdateProduct(w http.ResponseWriter, r *http.Request) {

	var err error

	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	claims, ok := utils.GetUserClaimsFromContext(r)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid user claims")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 50*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product id")
		return
	}

	var updatePrReq model.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&updatePrReq); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userEmail, ok := claims["email"].(string)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "User ID not found in token")
		return
	}

	curUser, err := c.usrSrvc.GetUserByEmail(ctx, userEmail)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "User not found")
		return
	}

	resID, err := c.prSrvc.UpdateProduct(ctx, updatePrReq, intId, curUser.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, resID)
}

func (c *MarketplaceController) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	var err error

	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	ctx, cancel := context.WithTimeout(r.Context(), 50*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["id"]
	intId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid id")
		return
	}

	err = c.prSrvc.DeleteProduct(ctx, intId)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}
