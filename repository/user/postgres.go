package repository

import (
	"context"
	"errors"

	"github.com/Leul-Michael/go-auth/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresUserRepo struct {
	db *gorm.DB
}

type UserRepo interface {
	GetById(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Insert(ctx context.Context, user model.User) error
	Update(ctx context.Context, user model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

var ErrNotExist = errors.New("user not found")

func (pr *PostgresUserRepo) GetById(ctx context.Context, id uuid.UUID) (model.User, error) {
	var user model.User

	err := pr.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Select("id, first_name, last_name, email, phone_number, last_login, is_deactivated").First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, ErrNotExist
		} else {
			return model.User{}, err
		}
	}

	return user, nil
}

func (pr *PostgresUserRepo) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	err := pr.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Select("id, first_name, last_name, email, phone_number, last_login, is_deactivated").First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.User{}, ErrNotExist
		} else {
			return model.User{}, err
		}
	}

	return user, nil
}

func (pr *PostgresUserRepo) Insert(ctx context.Context, user model.User) error {
	if err := pr.db.WithContext(ctx).Model(&model.User{}).Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (pr *PostgresUserRepo) Update(ctx context.Context, user model.User) error {
	if err := pr.db.WithContext(ctx).Model(&model.User{}).Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (pr *PostgresUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := pr.db.WithContext(ctx).Delete(&model.User{}, id).Error; err != nil {
		return err
	}
	return nil
}
