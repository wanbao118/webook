package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicateEmail = errors.New("email already exists")
	ErrRecordNotFound = errors.New("record not found")
)

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	Ctime    time.Time
	Utime    time.Time
}

// User represents the user model
type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (dao *UserDao) Insert(ctx context.Context, user *User) error {
	now := time.Now()
	user.Ctime = now
	user.Utime = now
	err := dao.db.WithContext(ctx).Create(user).Error

	if me, ok := err.(*mysql.MySQLError); ok {
		if me.Number == 1062 {
			return ErrDuplicateEmail
		}
	}

	return err

}

func (dao *UserDao) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}

	return &user, err
}

func (dao *UserDao) Update(ctx context.Context, user *User) error {
	err := dao.db.WithContext(ctx).Updates(user).Error
	return err
}

func (dao *UserDao) GetById(ctx context.Context, id int64) (*User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRecordNotFound
	}

	return &user, err
}
