package models

import (
	"time"

	"gorm.io/gorm"
)


type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);not null"`
	Email string `gorm:"unique"`
	Password string `gorm:"not null"`
	Role string `gorm:"type:varchar(255); not null"`
	Provider string `gorm:"not null"`
	Photo string `gorm:"not null"`
	VerificationCode string
	Verified bool `gorm:"not null"`
}

type SignUpInput struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
	Photo string `json:"photo" binding:"required"`
}

type SignInInput struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID int `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Role string `json:"role,omitempty"`
	Photo string `json:"photo,omitempty"`
	Provider string `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
