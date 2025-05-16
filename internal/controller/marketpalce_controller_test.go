package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vvwind/2025-MAI-Backend-V-Vetrov/internal/model"
	"github.com/vvwind/2025-MAI-Backend-V-Vetrov/internal/service"
)

func TestCreateUser(t *testing.T) {
	mockUserService := service.NewMockUserService(t)
	mockProductService := service.NewMockProductService(t)
	controller := NewMarketplaceController(mockProductService, mockUserService)

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "Success - create customer",
			requestBody: `{
				"name": "testuser",
				"email": "test@example.com",
				"password": "password123",
				"role": "customer"
			}`,
			mockSetup: func() {
				mockUserService.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
					Return(int64(14567), nil).Once()
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid role",
			requestBody: `{
				"name": "testuser",
				"email": "test@example.com",
				"password": "password123",
				"role": "invalid"
			}`,
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing fields",
			requestBody: `{
				"name": "",
				"email": "test@example.com",
				"password": "password123",
				"role": "customer"
			}`,
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Service error",
			requestBody: `{
				"name": "elon_petrov",
				"email": "elon_petrov@example.com",
				"password": "password123",
				"role": "customer"
			}`,
			mockSetup: func() {
				mockUserService.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
					Return(int64(0), errors.New("service error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest("POST", "/user", bytes.NewBufferString(tt.requestBody))
			rr := httptest.NewRecorder()
			controller.CreateUser(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockUserService.AssertExpectations(t)
		})
	}
}

func TestLoginUser(t *testing.T) {
	mockUserService := service.NewMockUserService(t)
	mockProductService := service.NewMockProductService(t)
	controller := NewMarketplaceController(mockProductService, mockUserService)

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "Success - valid login",
			requestBody: `{
				"email": "test@example.com",
				"password": "password123"
			}`,
			mockSetup: func() {
				mockUserService.On("LoginUser", mock.Anything, mock.Anything).
					Return("valid-token", nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid credentials",
			requestBody: `{
				"email": "test@example.com",
				"password": "wrongpassword"
			}`,
			mockSetup: func() {
				mockUserService.On("LoginUser", mock.Anything, mock.Anything).
					Return("", errors.New("invalid credentials")).Once()
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Missing fields",
			requestBody: `{
				"email": "",
				"password": "password123"
			}`,
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest("POST", "/user/login", bytes.NewBufferString(tt.requestBody))
			rr := httptest.NewRecorder()
			controller.LoginUser(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockUserService.AssertExpectations(t)
		})
	}
}

func TestGetProductByID(t *testing.T) {
	mockProductService := service.NewMockProductService(t)
	mockUserService := service.NewMockUserService(t)
	controller := NewMarketplaceController(mockProductService, mockUserService)

	tests := []struct {
		name           string
		productID      string
		mockSetup      func(id int64)
		expectedStatus int
	}{
		{
			name:      "Success - product found",
			productID: "1",
			mockSetup: func(id int64) {
				mockProductService.On("GetProductByID", mock.Anything, id).
					Return(&model.Product{ID: id, Title: "Test Product"}, nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "Product not found",
			productID: "999",
			mockSetup: func(id int64) {
				mockProductService.On("GetProductByID", mock.Anything, id).
					Return(nil, nil).Once()
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid ID",
			productID:      "invalid",
			mockSetup:      func(id int64) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "Service error",
			productID: "1",
			mockSetup: func(id int64) {
				mockProductService.On("GetProductByID", mock.Anything, id).
					Return(nil, errors.New("service error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				id, _ := strconv.ParseInt(tt.productID, 10, 64)
				tt.mockSetup(id)
			}

			req := httptest.NewRequest("GET", "/products/"+tt.productID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.productID})
			rr := httptest.NewRecorder()
			controller.GetProductByID(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockProductService.AssertExpectations(t)
		})
	}
}

func TestCreateProduct(t *testing.T) {
	mockProductService := service.NewMockProductService(t)
	mockUserService := service.NewMockUserService(t)
	controller := NewMarketplaceController(mockProductService, mockUserService)

	testName := strings.Replace(faker.Name(), " ", "", -1)
	testDomain := faker.DomainName()

	testEmail := fmt.Sprintf("%s@%s", testName, testDomain)
	testSeller := &model.User{ID: 1, Email: testEmail, Role: "seller"}

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name: "Success - valid product creation",
			requestBody: `{
                "title": "New Product",
                "product_description": "Description",
                "product_image": "image.jpg",
                "price": 100,
                "amount": 10
            }`,
			mockSetup: func() {
				mockUserService.On("GetUserByEmail", mock.Anything, testEmail).
					Return(testSeller, nil).Once()
				mockProductService.On("CreateProduct", mock.Anything, mock.Anything, *testSeller).
					Return(int64(1), nil).Once()
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Fail - invalid JSON",
			requestBody:    `{ invalid json }`,
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest("POST", "/products", bytes.NewBufferString(tt.requestBody))

			// Set up auth context
			claims := jwt.MapClaims{"email": testEmail}
			ctx := context.WithValue(req.Context(), "userClaims", claims)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			controller.CreateProduct(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockProductService.AssertExpectations(t)
			mockUserService.AssertExpectations(t)
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	mockProductService := service.NewMockProductService(t)
	mockUserService := service.NewMockUserService(t)
	controller := NewMarketplaceController(mockProductService, mockUserService)

	testName := strings.Replace(faker.Name(), " ", "", -1)
	testDomain := faker.DomainName()

	testEmail := fmt.Sprintf("%s@%s", testName, testDomain)
	testSeller := &model.User{ID: 1, Email: testEmail, Role: "seller", UserName: "elon_musk"}

	validUpdate := `{
        "title": "Updated Product",
        "product_description": "Updated Description",
        "product_image": "updated_pic.jpg",
        "price": 200,
        "amount": 20
    }`

	tests := []struct {
		name           string
		productID      string
		requestBody    string
		setupMocks     func()
		setupRequest   func(req *http.Request)
		expectedStatus int
	}{
		{
			name:        "Success - update product",
			productID:   "1",
			requestBody: validUpdate,
			setupMocks: func() {
				mockUserService.On("GetUserByEmail", mock.Anything, testEmail).
					Return(testSeller, nil).Once()
				mockProductService.On("UpdateProduct", mock.Anything, mock.Anything, int64(1), testSeller.ID).
					Return(int64(1), nil).Once()
			},
			setupRequest: func(req *http.Request) {
				claims := jwt.MapClaims{"email": testEmail}
				ctx := context.WithValue(req.Context(), "userClaims", claims)
				*req = *req.WithContext(ctx)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "Invalid product ID",
			productID:   "invalid",
			requestBody: validUpdate,
			setupMocks: func() {
			},
			setupRequest: func(req *http.Request) {
				claims := jwt.MapClaims{"email": testEmail}
				ctx := context.WithValue(req.Context(), "userClaims", claims)
				*req = *req.WithContext(ctx)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Invalid update data",
			productID:   "2",
			requestBody: `{"title": ""}`,
			setupMocks: func() {
			},
			setupRequest: func(req *http.Request) {
				claims := jwt.MapClaims{"email": testEmail}
				ctx := context.WithValue(req.Context(), "userClaims", claims)
				*req = *req.WithContext(ctx)
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Unauthorized - no auth context",
			productID:      "3",
			requestBody:    validUpdate,
			setupMocks:     func() {},
			setupRequest:   func(req *http.Request) {},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.setupMocks()

			// Create request
			req := httptest.NewRequest("PUT", "/products/"+tt.productID, bytes.NewBufferString(tt.requestBody))
			req = mux.SetURLVars(req, map[string]string{"id": tt.productID})

			// Setup request context
			tt.setupRequest(req)

			rr := httptest.NewRecorder()
			controller.UpdateProduct(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus != http.StatusUnauthorized && tt.expectedStatus != http.StatusBadRequest {
				mockProductService.AssertExpectations(t)
			}
			mockUserService.AssertExpectations(t)
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	mockProductService := service.NewMockProductService(t)
	mockUserService := service.NewMockUserService(t)
	controller := NewMarketplaceController(mockProductService, mockUserService)

	tests := []struct {
		name           string
		productID      string
		mockSetup      func(id int64)
		expectedStatus int
	}{
		{
			name:      "Success - delete product",
			productID: "123",
			mockSetup: func(id int64) {
				mockProductService.On("DeleteProduct", mock.Anything, id).
					Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid ID",
			productID:      "invalid",
			mockSetup:      func(id int64) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "Service error",
			productID: "1",
			mockSetup: func(id int64) {
				mockProductService.On("DeleteProduct", mock.Anything, id).
					Return(errors.New("service error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				id, _ := strconv.ParseInt(tt.productID, 10, 64)
				tt.mockSetup(id)
			}

			req := httptest.NewRequest("DELETE", "/products/"+tt.productID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.productID})
			rr := httptest.NewRecorder()
			controller.DeleteProduct(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockProductService.AssertExpectations(t)
		})
	}
}

func TestAddToCart(t *testing.T) {
	mockProductService := service.NewMockProductService(t)
	mockUserService := service.NewMockUserService(t)
	controller := NewMarketplaceController(mockProductService, mockUserService)

	testName := strings.Replace(faker.Name(), " ", "", -1)
	testDomain := faker.DomainName()

	testEmail := fmt.Sprintf("%s@%s", testName, testDomain)
	testCustomer := &model.User{ID: 1, Email: testEmail, Role: "customer"}

	tests := []struct {
		name           string
		productID      string
		mockSetup      func(productID, userID int64)
		expectedStatus int
	}{
		{
			name:      "Success - add to cart",
			productID: "1",
			mockSetup: func(productID, userID int64) {
				mockUserService.On("GetUserByEmail", mock.Anything, testEmail).
					Return(testCustomer, nil).Once()
				mockProductService.On("AddToCart", mock.Anything, productID, userID).
					Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid product ID",
			productID:      "invalid",
			mockSetup:      func(productID, userID int64) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "Service error",
			productID: "2",
			mockSetup: func(productID, userID int64) {
				mockUserService.On("GetUserByEmail", mock.Anything, testEmail).
					Return(testCustomer, nil).Once()
				mockProductService.On("AddToCart", mock.Anything, productID, userID).
					Return(errors.New("service error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				productID, _ := strconv.ParseInt(tt.productID, 10, 64)
				tt.mockSetup(productID, testCustomer.ID)
			}

			req := httptest.NewRequest("POST", "/products/cart/"+tt.productID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.productID})

			// Set up auth context
			claims := jwt.MapClaims{"email": testEmail}
			ctx := context.WithValue(req.Context(), "userClaims", claims)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			controller.AddToCart(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockProductService.AssertExpectations(t)
			mockUserService.AssertExpectations(t)
		})
	}
}

func TestBuyProduct(t *testing.T) {
	mockProductService := service.NewMockProductService(t)
	mockUserService := service.NewMockUserService(t)
	controller := NewMarketplaceController(mockProductService, mockUserService)

	testName := strings.Replace(faker.Name(), " ", "", -1)
	testDomain := faker.DomainName()

	testEmail := fmt.Sprintf("%s@%s", testName, testDomain)
	testCustomer := &model.User{
		ID:       123,
		Email:    testEmail,
		Role:     "customer",
		UserName: "ivan_ivanov",
	}

	tests := []struct {
		name           string
		productID      string
		mockSetup      func(productID, userID int64)
		expectedStatus int
	}{
		{
			name:      "Success - buy product",
			productID: "1",
			mockSetup: func(productID, userID int64) {
				mockUserService.On("GetUserByEmail", mock.Anything, testEmail).
					Return(testCustomer, nil).Once()
				mockProductService.On("BuyProduct", mock.Anything, productID, userID).
					Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid product ID",
			productID:      "invalid",
			mockSetup:      func(productID, userID int64) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "Service error",
			productID: "2",
			mockSetup: func(productID, userID int64) {
				mockUserService.On("GetUserByEmail", mock.Anything, testEmail).
					Return(testCustomer, nil).Once()
				mockProductService.On("BuyProduct", mock.Anything, productID, userID).
					Return(errors.New("service error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				productID, _ := strconv.ParseInt(tt.productID, 10, 64)
				tt.mockSetup(productID, testCustomer.ID)
			}

			req := httptest.NewRequest("POST", "/products/buy/"+tt.productID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.productID})

			// Set up auth context
			claims := jwt.MapClaims{"email": testEmail}
			ctx := context.WithValue(req.Context(), "userClaims", claims)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			controller.BuyProduct(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			mockProductService.AssertExpectations(t)
			mockUserService.AssertExpectations(t)
		})
	}
}
