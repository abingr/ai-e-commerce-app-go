package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"ai-e-commerce-app-go/backend/internal/models"
	"ai-e-commerce-app-go/backend/internal/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type UserRepository interface {
	Create(ctx context.Context, user models.User) (models.User, error)
	FindByEmail(ctx context.Context, email string) (models.User, error)
	FindByID(ctx context.Context, id string) (models.User, error)
}

type AuthClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct {
	repository UserRepository
	jwtSecret  []byte
	jwtIssuer  string
}

func NewAuthService(repository UserRepository, jwtSecret string, jwtIssuer string) AuthService {
	return AuthService{
		repository: repository,
		jwtSecret:  []byte(jwtSecret),
		jwtIssuer:  jwtIssuer,
	}
}

func (s AuthService) Register(ctx context.Context, input models.RegisterUserInput) (models.AuthResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.AuthResponse{}, err
	}

	user := models.User{
		Name:         strings.TrimSpace(input.Name),
		Email:        normalizeEmail(input.Email),
		PasswordHash: string(passwordHash),
		Role:         models.RoleCustomer,
	}

	createdUser, err := s.repository.Create(ctx, user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	token, err := s.GenerateToken(createdUser)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return models.AuthResponse{
		User:  createdUser,
		Token: token,
	}, nil
}

func (s AuthService) Login(ctx context.Context, input models.LoginUserInput) (models.AuthResponse, error) {
	user, err := s.repository.FindByEmail(ctx, normalizeEmail(input.Email))
	if errors.Is(err, repositories.ErrNotFound) {
		return models.AuthResponse{}, ErrInvalidCredentials
	}

	if err != nil {
		return models.AuthResponse{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return models.AuthResponse{}, ErrInvalidCredentials
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return models.AuthResponse{}, err
	}

	return models.AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s AuthService) GetUser(ctx context.Context, id string) (models.User, error) {
	return s.repository.FindByID(ctx, id)
}

func (s AuthService) GenerateToken(user models.User) (string, error) {
	now := time.Now()
	claims := AuthClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.jwtIssuer,
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s AuthService) ParseToken(tokenString string) (AuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}

		return s.jwtSecret, nil
	}, jwt.WithIssuer(s.jwtIssuer))
	if err != nil {
		return AuthClaims{}, err
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		return AuthClaims{}, errors.New("invalid token")
	}

	return *claims, nil
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
