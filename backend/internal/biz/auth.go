package biz

import (
	"context"
	"errors"

	"quant-trader/api/middleware"
	"quant-trader/internal/model"
	"quant-trader/internal/repo"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type Auth struct {
	user   *repo.User
	logger *zap.Logger
}

func NewAuth(user *repo.User, logger *zap.Logger) *Auth {
	return &Auth{
		user:   user,
		logger: logger,
	}
}

func (b *Auth) Register(ctx context.Context, email, password string) (*model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := b.user.Create(ctx, user); err != nil {
		b.logger.Error("failed to register user", zap.Error(err))
		return nil, ErrEmailAlreadyExists
	}

	return user, nil
}

func (b *Auth) Login(ctx context.Context, email, password string) (string, error) {
	user, err := b.user.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (b *Auth) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	return b.user.GetByID(ctx, id)
}
