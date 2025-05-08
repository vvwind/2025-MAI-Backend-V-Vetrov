package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"rest-api/internal/model"
	"rest-api/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	LoginUser(ctx context.Context, usr model.UserLogin) (string, error)
	CreateUser(ctx context.Context, usr model.UserRegister, hashedPassword string) (int64, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

type JWTClaims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func (s *userService) CreateUser(ctx context.Context, usr model.UserRegister, hashedPassword string) (int64, error) {
	return s.userRepo.CreateUser(ctx, FromRequestToModel(usr), hashedPassword)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}

func (s *userService) LoginUser(ctx context.Context, login model.UserLogin) (string, error) {
	// 1. Get hashed password from repository
	hashedPassword, err := s.userRepo.GetHashedPassword(ctx, login.Email)
	if err != nil {
		return "", fmt.Errorf("failed to get user credentials: %w", err)
	}

	// 2. Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(login.Password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 3. Get full user details
	user, err := s.userRepo.GetUserByEmail(ctx, login.Email)
	if err != nil {
		return "", fmt.Errorf("failed to get user details: %w", err)
	}

	// 4. Generate JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *userService) generateJWT(user *model.User) (string, error) {
	// Create JWT claims
	claims := JWTClaims{
		UserID:   user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "HomeBerries",
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Get secret key from environment
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET not set in environment")
	}

	// Sign and get the complete encoded token as string
	return token.SignedString([]byte(secret))
}
