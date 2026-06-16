package models

import "time"

// User mirrors the users table.
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture,omitempty"`
	Language string `json:"language"`
}

// UserBudget mirrors the user_budgets table.
type UserBudget struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Note       *string `json:"note,omitempty"`
	Icon       string  `json:"icon"`
	Color      string  `json:"color"`
	Percentage float64 `json:"percentage"`
	Amount     float64 `json:"amount"`
}

// BudgetBooking mirrors the budget_bookings table.
type BudgetBooking struct {
	ID         int64     `json:"id"`
	BudgetID   *int64    `json:"budget_id"`
	BudgetName *string   `json:"budget_name"`
	IsIncome   bool      `json:"is_income"`
	Amount     float64   `json:"amount"`
	Note       *string   `json:"note,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
