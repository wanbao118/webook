package repository

import (
	"context"

	"github.com/webook/internal/domain"
	"github.com/webook/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{dao: dao}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.dao.Insert(ctx, &dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := r.dao.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return r.toDomain(user), nil
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (*domain.User, error) {
	user, err := r.dao.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toDomain(user), nil
}

func (r *UserRepository) Update(ctx context.Context, user domain.User) error {
	return r.dao.Update(ctx, &dao.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	})
}

func (r *UserRepository) toDomain(user *dao.User) *domain.User {
	return &domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}
}
