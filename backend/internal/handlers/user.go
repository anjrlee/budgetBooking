package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"budgetbooking/internal/googleauth"
	"budgetbooking/internal/models"
)

const (
	defaultLanguage    = "zh-Hant"
	savingBudgetIcon   = "savings"
	savingBudgetColor  = "#B5EAD7"
	savingBudgetPctMax = 100
)

// Login handles POST /user/login.
func (h *Handlers) Login(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing google id token"})
		return
	}

	ctx := c.Request.Context()
	claims, err := googleauth.Verify(ctx, req.Token, h.GoogleClientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid google id token"})
		return
	}

	tx, err := h.DB.Begin(ctx)
	if err != nil {
		respondError(c, err)
		return
	}
	defer tx.Rollback(ctx)

	user, budgets, err := findOrCreateUser(ctx, tx, claims)
	if err != nil {
		respondError(c, err)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		respondError(c, err)
		return
	}

	token, err := h.issueJWT(user.ID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"user":    user,
		"budgets": budgets,
	})
}

// findOrCreateUser looks up the user by Google subject, creating the user
// record and default "saving" budget on first login.
func findOrCreateUser(ctx context.Context, tx pgx.Tx, claims *googleauth.Claims) (models.User, []models.UserBudget, error) {
	var (
		user    models.User
		picture *string
	)
	err := tx.QueryRow(ctx, `
		SELECT id, email, name, picture, language
		FROM users WHERE id = $1`, claims.Subject).
		Scan(&user.ID, &user.Email, &user.Name, &picture, &user.Language)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		user = models.User{
			ID:       claims.Subject,
			Email:    claims.Email,
			Name:     claims.Name,
			Picture:  claims.Picture,
			Language: defaultLanguage,
		}
		var picturePtr *string
		if user.Picture != "" {
			picturePtr = &user.Picture
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO users (id, email, name, picture, language)
			VALUES ($1, $2, $3, $4, $5)`,
			user.ID, user.Email, user.Name, picturePtr, user.Language); err != nil {
			return models.User{}, nil, err
		}

		savingID, err := insertBudget(ctx, tx, user.ID, savingBudgetName, nil, savingBudgetIcon, savingBudgetColor, savingBudgetPctMax)
		if err != nil {
			return models.User{}, nil, err
		}

		budgets := []models.UserBudget{{
			ID:         savingID,
			Name:       savingBudgetName,
			Icon:       savingBudgetIcon,
			Color:      savingBudgetColor,
			Percentage: savingBudgetPctMax,
			Amount:     0,
		}}
		return user, budgets, nil

	case err != nil:
		return models.User{}, nil, err

	default:
		if picture != nil {
			user.Picture = *picture
		}
		budgets, err := getBudgets(ctx, tx, user.ID)
		if err != nil {
			return models.User{}, nil, err
		}
		return user, budgets, nil
	}
}

// Me handles GET /user/me.
func (h *Handlers) Me(c *gin.Context) {
	ctx := c.Request.Context()
	userID := currentUserID(c)

	var (
		user    models.User
		picture *string
	)
	err := h.DB.QueryRow(ctx, `
		SELECT id, email, name, picture, language
		FROM users WHERE id = $1`, userID).
		Scan(&user.ID, &user.Email, &user.Name, &picture, &user.Language)
	if err != nil {
		respondError(c, err)
		return
	}
	if picture != nil {
		user.Picture = *picture
	}

	budgets, err := getBudgets(ctx, h.DB, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"budgets": budgets,
	})
}

// UpdateSetting handles POST /user/setting.
func (h *Handlers) UpdateSetting(c *gin.Context) {
	var req struct {
		Language string `json:"language" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "missing language"})
		return
	}
	if req.Language != "zh-Hant" && req.Language != "en" {
		respondError(c, badRequest("unsupported language %q", req.Language))
		return
	}

	ctx := c.Request.Context()
	userID := currentUserID(c)
	if _, err := h.DB.Exec(ctx, `UPDATE users SET language = $1 WHERE id = $2`, req.Language, userID); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"language": req.Language})
}
