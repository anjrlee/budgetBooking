package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddBudget handles POST /setBudget/add.
func (h *Handlers) AddBudget(c *gin.Context) {
	userID := currentUserID(c)

	var req struct {
		Name               string             `json:"name" binding:"required"`
		Note               *string            `json:"note"`
		Icon               string             `json:"icon" binding:"required"`
		Color              string             `json:"color" binding:"required"`
		Percentage         float64            `json:"percentage"`
		UpdatedPercentages map[string]float64 `json:"updatedPercentages"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}
	if req.Name == savingBudgetName {
		respondError(c, badRequest("%q is a reserved budget name", savingBudgetName))
		return
	}

	ctx := c.Request.Context()
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		respondError(c, err)
		return
	}
	defer tx.Rollback(ctx)

	newID, err := insertBudget(ctx, tx, userID, req.Name, req.Note, req.Icon, req.Color, req.Percentage)
	if err != nil {
		if isUniqueViolation(err) {
			respondError(c, badRequest("a budget with that name, icon, or color already exists"))
			return
		}
		respondError(c, err)
		return
	}

	updated := make(map[string]float64, len(req.UpdatedPercentages)+1)
	for k, v := range req.UpdatedPercentages {
		updated[k] = v
	}
	updated[strconv.FormatInt(newID, 10)] = req.Percentage

	if err := syncPercentages(ctx, tx, userID, updated); err != nil {
		respondError(c, err)
		return
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

// ChangeBudget handles POST /setBudget/change.
func (h *Handlers) ChangeBudget(c *gin.Context) {
	userID := currentUserID(c)

	var req struct {
		ID                 int64              `json:"id" binding:"required"`
		Name               string             `json:"name" binding:"required"`
		Note               *string            `json:"note"`
		Icon               string             `json:"icon" binding:"required"`
		Color              string             `json:"color" binding:"required"`
		UpdatedPercentages map[string]float64 `json:"updatedPercentages"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}
	if req.Name == savingBudgetName {
		respondError(c, badRequest("%q is a reserved budget name", savingBudgetName))
		return
	}

	ctx := c.Request.Context()
	tx, err := h.DB.Begin(ctx)
	if err != nil {
		respondError(c, err)
		return
	}
	defer tx.Rollback(ctx)

	target, err := getBudget(ctx, tx, userID, req.ID)
	if err != nil {
		respondError(c, err)
		return
	}
	if target.Name == savingBudgetName {
		respondError(c, badRequest("the saving bucket cannot be modified"))
		return
	}

	if _, err := tx.Exec(ctx, `
		UPDATE user_budgets SET name = $1, note = $2, icon = $3, color = $4
		WHERE id = $5 AND user_id = $6`, req.Name, req.Note, req.Icon, req.Color, req.ID, userID); err != nil {
		if isUniqueViolation(err) {
			respondError(c, badRequest("a budget with that name, icon, or color already exists"))
			return
		}
		respondError(c, err)
		return
	}

	if len(req.UpdatedPercentages) > 0 {
		if err := syncPercentages(ctx, tx, userID, req.UpdatedPercentages); err != nil {
			respondError(c, err)
			return
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

// DeleteBudget handles POST /setBudget/delete.
func (h *Handlers) DeleteBudget(c *gin.Context) {
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

	target, err := getBudget(ctx, tx, userID, req.ID)
	if err != nil {
		respondError(c, err)
		return
	}
	if target.Name == savingBudgetName {
		respondError(c, badRequest("the saving bucket cannot be deleted"))
		return
	}

	if _, err := tx.Exec(ctx, `
		UPDATE user_budgets SET amount = amount + $1, percentage = percentage + $2
		WHERE user_id = $3 AND name = $4`, target.Amount, target.Percentage, userID, savingBudgetName); err != nil {
		respondError(c, err)
		return
	}

	if _, err := tx.Exec(ctx, `DELETE FROM user_budgets WHERE id = $1 AND user_id = $2`, req.ID, userID); err != nil {
		respondError(c, err)
		return
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
