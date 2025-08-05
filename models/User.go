package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName      string    `gorm:"size:255;not null" json:"first_name" validate:"required,min=2,max=100"`
	LastName       string    `gorm:"size:255;not null" json:"last_name" validate:"required,min=2,max=100"`
	Password       string    `gorm:"size:255;not null" json:"password" validate:"required,min=8"`
	Email          string    `gorm:"size:255;not null;unique" json:"email" validate:"required,email"`
	Phone          string    `gorm:"size:255" json:"phone" validate:"required"`
	Token          string    `gorm:"size:255" json:"token"`
	Role           string    `gorm:"size:50;not null;default:'user'" json:"role" validate:"required,oneof=admin user"`
	RefreshToken   string    `gorm:"size:255" json:"refresh_token"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	UserID         string    `gorm:"size:255;not null;unique" json:"user_id"`
}