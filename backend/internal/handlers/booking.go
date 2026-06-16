package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"budgetbooking/internal/models"
)

// AddBooking handles POST /booking/add.
func (h *Handlers) AddBooking(c *gin.Context) {
	userID := currentUserID(c)

	var req struct {
		IsIncome bool    `json:"is_income"`
		Amount   float64 `json:"amount"`
		Note     *string `json:"note"`
		BudgetID *int64  `json:"budget_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}
	if req.Amount <= 0 {
		respondError(c, badRequest("amount must be positive"))
		return
	}

	ctx := c.Request.Context()
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		respondError(c, err)
		return
	}
	defer tx.Rollback(ctx)

	switch {
	case !req.IsIncome:
		if req.BudgetID == nil {
			respondError(c, badRequest("budget_id is required for spending"))
			return
		}
		budget, err := getBudget(ctx, tx, userID, *req.BudgetID)
		if err != nil {
			respondError(c, err)
			return
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO budget_bookings (user_id, budget_id, budget_name, is_income, amount, note)
			VALUES ($1, $2, $3, false, $4, $5)`, userID, budget.ID, budget.Name, req.Amount, req.Note); err != nil {
			respondError(c, err)
			return
		}
		if _, err := tx.Exec(ctx, `UPDATE user_budgets SET amount = amount - $1 WHERE id = $2`, req.Amount, budget.ID); err != nil {
			respondError(c, err)
			return
		}

	case req.BudgetID != nil:
		budget, err := getBudget(ctx, tx, userID, *req.BudgetID)
		if err != nil {
			respondError(c, err)
			return
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO budget_bookings (user_id, budget_id, budget_name, is_income, amount, note)
			VALUES ($1, $2, $3, true, $4, $5)`, userID, budget.ID, budget.Name, req.Amount, req.Note); err != nil {
			respondError(c, err)
			return
		}
		if _, err := tx.Exec(ctx, `UPDATE user_budgets SET amount = amount + $1 WHERE id = $2`, req.Amount, budget.ID); err != nil {
			respondError(c, err)
			return
		}

	default:
		if _, err := tx.Exec(ctx, `
			INSERT INTO budget_bookings (user_id, budget_id, budget_name, is_income, amount, note)
			VALUES ($1, NULL, NULL, true, $2, $3)`, userID, req.Amount, req.Note); err != nil {
			respondError(c, err)
			return
		}

		budgets, err := getBudgets(ctx, tx, userID)
		if err != nil {
			respondError(c, err)
			return
		}
		for id, alloc := range distributeAmounts(budgets, req.Amount) {
			if _, err := tx.Exec(ctx, `UPDATE user_budgets SET amount = amount + $1 WHERE id = $2`, alloc, id); err != nil {
				respondError(c, err)
				return
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		respondError(c, err)
		return
	}

	budgets, err := getBudgets(ctx, h.DB, userID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"budgets": budgets})
}

// DeleteBooking handles POST /booking/delete.
func (h *Handlers) DeleteBooking(c *gin.Context) {
	userID := currentUserID(c)

	var req struct {
		ID int64 `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	ctx := c.Request.Context()
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		respondError(c, err)
		return
	}
	defer tx.Rollback(ctx)

	var booking models.BudgetBooking
	err = tx.QueryRow(ctx, `
		SELECT id, budget_id, budget_name, is_income, amount
		FROM budget_bookings WHERE id = $1 AND user_id = $2`, req.ID, userID).
		Scan(&booking.ID, &booking.BudgetID, &booking.BudgetName, &booking.IsIncome, &booking.Amount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondError(c, notFound("booking not found"))
			return
		}
		respondError(c, err)
		return
	}

	if _, err := tx.Exec(ctx, `DELETE FROM budget_bookings WHERE id = $1 AND user_id = $2`, req.ID, userID); err != nil {
		respondError(c, err)
		return
	}

	switch {
	case !booking.IsIncome && booking.BudgetID != nil:
		// Spending on an existing budget: refund it.
		if _, err := tx.Exec(ctx, `UPDATE user_budgets SET amount = amount + $1 WHERE id = $2`, booking.Amount, *booking.BudgetID); err != nil {
			respondError(c, err)
			return
		}

	case !booking.IsIncome:
		// Spending on a since-deleted budget: refund went to saving.
		saving, err := getSavingBudget(ctx, tx, userID)
		if err != nil {
			respondError(c, err)
			return
		}
		if _, err := tx.Exec(ctx, `UPDATE user_budgets SET amount = amount + $1 WHERE id = $2`, booking.Amount, saving.ID); err != nil {
			respondError(c, err)
			return
		}

	case booking.IsIncome && booking.BudgetID != nil:
		// Direct income to an existing budget: deduct it.
		if _, err := tx.Exec(ctx, `UPDATE user_budgets SET amount = amount - $1 WHERE id = $2`, booking.Amount, *booking.BudgetID); err != nil {
			respondError(c, err)
			return
		}

	case booking.IsIncome && booking.BudgetName != nil:
		// Direct income to a since-deleted budget: it landed in saving.
		saving, err := getSavingBudget(ctx, tx, userID)
		if err != nil {
			respondError(c, err)
			return
		}
		if _, err := tx.Exec(ctx, `UPDATE user_budgets SET amount = amount - $1 WHERE id = $2`, booking.Amount, saving.ID); err != nil {
			respondError(c, err)
			return
		}

	default:
		// Distributed income: reverse using current percentages.
		budgets, err := getBudgets(ctx, tx, userID)
		if err != nil {
			respondError(c, err)
			return
		}
		for id, alloc := range distributeAmounts(budgets, booking.Amount) {
			if _, err := tx.Exec(ctx, `UPDATE user_budgets SET amount = amount - $1 WHERE id = $2`, alloc, id); err != nil {
				respondError(c, err)
				return
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		respondError(c, err)
		return
	}

	budgets, err := getBudgets(ctx, h.DB, userID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"budgets": budgets})
}
