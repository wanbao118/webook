package service

import (
	"context"
	"errors"

	"github.com/webook/internal/domain"
	"github.com/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("invalid user or password")
)

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repository: repo,
	}
}

func (svc *UserService) Signup(ctx context.Context, user domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return svc.repository.Create(ctx, &user)
}

func (svc *UserService) Signin(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := svc.repository.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return &domain.User{}, errors.New("user not found")
	}
	if err != nil {
		return &domain.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return &domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil
}
