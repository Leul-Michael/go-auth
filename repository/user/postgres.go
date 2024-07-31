package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Leul-Michael/go-auth/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresUserRepo struct {
	DB *gorm.DB
}

type UserRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (*model.User, error)
	EmailExists(ctx context.Context, email string) int64
	ComparePassword(ctx context.Context, sub interface{}) (*Sub, error)
	Insert(ctx context.Context, user model.User) error
	Update(ctx context.Context, user model.User) error
	UpdateField(ctx context.Context, id uuid.UUID, field string, value interface{}) error
	Delete(ctx context.Context, id uuid.UUID) error
}

var ErrNotExist = errors.New("user not found")

func (pr *PostgresUserRepo) GetById(ctx context.Context, id uuid.UUID) (model.User, error) {
	var user model.User

	err := pr.DB.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Select("id, first_name, last_name, email, phone_number, last_login, is_deactivated").First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, ErrNotExist
		} else {
			return model.User{}, fmt.Errorf("failed to get record: %w", err)
		}
	}

	return user, nil
}

func (pr *PostgresUserRepo) EmailExists(ctx context.Context, email string) int64 {
	var count int64
	pr.DB.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count)
	return count
}

type Sub struct {
	Id            uuid.UUID `json:"id"`
	Password      string    `json:"password"`
	IsDeactivated bool      `json:"is_deactivated"`
}

func (pr *PostgresUserRepo) ComparePassword(ctx context.Context, email string) (*Sub, error) {
	var sub Sub
	err := pr.DB.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Select("id, password, is_deactivated").Scan(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (pr *PostgresUserRepo) Insert(ctx context.Context, user model.User) error {
	if err := pr.DB.WithContext(ctx).Model(&model.User{}).Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (pr *PostgresUserRepo) Update(ctx context.Context, user model.User) error {
	if err := pr.DB.WithContext(ctx).Model(&model.User{}).Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (pr *PostgresUserRepo) UpdateField(ctx context.Context, id uuid.UUID, field string, value interface{}) error {
	if err := pr.DB.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update(field, value).
		Error; err != nil {
		return err
	}
	return nil
}

func (pr *PostgresUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := pr.DB.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		return err
	}
	return nil
}
