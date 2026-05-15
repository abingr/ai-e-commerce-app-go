package services_test

import (
	"context"
	"errors"
	"testing"

	"ai-e-commerce-app-go/backend/internal/models"
	"ai-e-commerce-app-go/backend/internal/repositories"
	"ai-e-commerce-app-go/backend/internal/services"

	"golang.org/x/crypto/bcrypt"
)

type stubUserRepository struct {
	create      func(ctx context.Context, user models.User) (models.User, error)
	findByEmail func(ctx context.Context, email string) (models.User, error)
	findByID    func(ctx context.Context, id string) (models.User, error)
}

func (r stubUserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	return r.create(ctx, user)
}

func (r stubUserRepository) FindByEmail(ctx context.Context, email string) (models.User, error) {
	return r.findByEmail(ctx, email)
}

func (r stubUserRepository) FindByID(ctx context.Context, id string) (models.User, error) {
	return r.findByID(ctx, id)
}

func TestAuthServiceRegisterCreatesCustomerAndToken(t *testing.T) {
	repository := stubUserRepository{
		create: func(ctx context.Context, user models.User) (models.User, error) {
			if user.Email != "learner@example.com" {
				t.Fatalf("expected normalized email, got %q", user.Email)
			}

			if user.Role != models.RoleCustomer {
				t.Fatalf("expected customer role, got %q", user.Role)
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password123")); err != nil {
				t.Fatalf("expected password to be hashed correctly: %v", err)
			}

			user.ID = "f5e5a9b0-4055-421f-bd6f-7e755a16d1af"
			return user, nil
		},
	}

	service := services.NewAuthService(repository, "test-secret", "test-issuer")

	response, err := service.Register(context.Background(), models.RegisterUserInput{
		Name:     "Learner",
		Email:    " LEARNER@example.com ",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("register user: %v", err)
	}

	if response.Token == "" {
		t.Fatal("expected token")
	}

	claims, err := service.ParseToken(response.Token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}

	if claims.UserID != response.User.ID {
		t.Fatalf("expected token user id %q, got %q", response.User.ID, claims.UserID)
	}
}

func TestAuthServiceLoginRejectsWrongPassword(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	repository := stubUserRepository{
		findByEmail: func(ctx context.Context, email string) (models.User, error) {
			return models.User{
				ID:           "f5e5a9b0-4055-421f-bd6f-7e755a16d1af",
				Email:        email,
				PasswordHash: string(passwordHash),
				Role:         models.RoleCustomer,
			}, nil
		},
	}

	service := services.NewAuthService(repository, "test-secret", "test-issuer")

	_, err = service.Login(context.Background(), models.LoginUserInput{
		Email:    "learner@example.com",
		Password: "wrong-password",
	})
	if !errors.Is(err, services.ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestAuthServiceLoginRejectsUnknownEmail(t *testing.T) {
	repository := stubUserRepository{
		findByEmail: func(ctx context.Context, email string) (models.User, error) {
			return models.User{}, repositories.ErrNotFound
		},
	}

	service := services.NewAuthService(repository, "test-secret", "test-issuer")

	_, err := service.Login(context.Background(), models.LoginUserInput{
		Email:    "missing@example.com",
		Password: "password123",
	})
	if !errors.Is(err, services.ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestAuthServiceParseTokenRejectsWrongIssuer(t *testing.T) {
	repository := stubUserRepository{}
	service := services.NewAuthService(repository, "test-secret", "expected-issuer")
	otherIssuerService := services.NewAuthService(repository, "test-secret", "other-issuer")

	token, err := otherIssuerService.GenerateToken(models.User{
		ID:    "f5e5a9b0-4055-421f-bd6f-7e755a16d1af",
		Email: "learner@example.com",
		Role:  models.RoleCustomer,
	})
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	_, err = service.ParseToken(token)
	if err == nil {
		t.Fatal("expected wrong issuer token to be rejected")
	}
}
