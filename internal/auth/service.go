package auth

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"

	"github.com/jackc/pgx/v5"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service interface {
	Register(ctx context.Context, input RegisterInput) (Account, error)
	Login(ctx context.Context, input LoginInput) (Account, error)
	Unregister(ctx context.Context, input UnregisterInput) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(ctx context.Context, input RegisterInput) (Account, error) {
	hash := hashPassword(input.Password)
	return s.repo.Create(ctx, input.Username, hash)
}

func (s *service) Login(ctx context.Context, input LoginInput) (Account, error) {
	account, err := s.repo.GetByUsername(ctx, input.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Account{}, ErrInvalidCredentials
		}
		return Account{}, err
	}

	if !verifyPassword(input.Password, account.PasswordHash) {
		return Account{}, ErrInvalidCredentials
	}

	return account, nil
}

func (s *service) Unregister(ctx context.Context, input UnregisterInput) error {
	account, err := s.Login(ctx, LoginInput{Username: input.Username, Password: input.Password})
	if err != nil {
		return err
	}

	return s.repo.DeleteByID(ctx, account.ID)
}

func hashPassword(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func verifyPassword(password string, storedHash string) bool {
	calculated := hashPassword(password)
	return subtle.ConstantTimeCompare([]byte(calculated), []byte(storedHash)) == 1
}
