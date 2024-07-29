package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base definition for all models
type Base struct {
	ID        uuid.UUID      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleUser   Role = "user"
	RoleClient Role = "client"
)

type User struct {
	Base
	FirstName     string     `gorm:"size:100;not null" json:"first_name"`
	LastName      string     `gorm:"size:100;not null" json:"last_name"`
	Email         string     `gorm:"size:100;unique;not null"`
	Password      string     `gorm:"not null" json:"-"`
	PhoneNumber   *string    `json:"phone_number"`
	LastLogin     *time.Time `json:"last_login"`
	Role          Role       `gorm:"type:role;default:'user'" json:"role"`
	IsDeactivated bool       `gorm:"default:false" json:"is_deactivated"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
