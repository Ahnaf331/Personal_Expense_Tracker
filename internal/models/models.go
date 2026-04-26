package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string     `json:"name" gorm:"not null"`
	Email      string     `json:"email" gorm:"uniqueIndex;not null"`
	Password   string     `json:"-" gorm:"not null"`
	Expenses   []Expense  `json:"-"`
	Categories []Category `json:"-"`
}

type Category struct {
	gorm.Model
	UserID   uint      `json:"user_id"`
	Name     string    `json:"name" gorm:"not null"`
	Budget   float64   `json:"budget"`
	Expenses []Expense `json:"-"`
}

type Expense struct {
	gorm.Model
	UserID      uint      `json:"user_id"`
	CategoryID  uint      `json:"category_id"`
	Category    Category  `json:"category" gorm:"foreignKey:CategoryID"`
	Amount      float64   `json:"amount" gorm:"not null"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateExpenseRequest struct {
	CategoryID  uint    `json:"category_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
	Date        string  `json:"date" binding:"required"`
}

type UpdateExpenseRequest struct {
	CategoryID  uint    `json:"category_id"`
	Amount      float64 `json:"amount" binding:"omitempty,gt=0"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

type CreateCategoryRequest struct {
	Name   string  `json:"name" binding:"required"`
	Budget float64 `json:"budget"`
}

type UpdateCategoryRequest struct {
	Name   string  `json:"name"`
	Budget float64 `json:"budget"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
