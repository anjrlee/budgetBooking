package handlers

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"budgetbooking/internal/middleware"
	"budgetbooking/internal/models"
)

const savingBudgetName = "saving"

const percentageEpsilon = 0.001

// querier is satisfied by both *pgxpool.Pool and pgx.Tx, letting helpers run
// either standalone or inside a caller-managed transaction.
type querier interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

// apiError carries an HTTP status alongside a user-facing message.
type apiError struct {
	status  int
	message string
}

func (e *apiError) Error() string { return e.message }

func badRequest(format string, args ...interface{}) error {
	return &apiError{status: http.StatusBadRequest, message: fmt.Sprintf(format, args...)}
}

func notFound(format string, args ...interface{}) error {
	return &apiError{status: http.StatusNotFound, message: fmt.Sprintf(format, args...)}
}

func unauthorized(format string, args ...interface{}) error {
	return &apiError{status: http.StatusUnauthorized, message: fmt.Sprintf(format, args...)}
}

// respondError writes the appropriate JSON error response for err.
func respondError(c *gin.Context, err error) {
	var ae *apiError
	if errors.As(err, &ae) {
		c.JSON(ae.status, gin.H{"message": ae.message})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
}

// isUniqueViolation reports whether err is a Postgres unique constraint violation.
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// currentUserID extracts the authenticated user's ID from the gin context.
func currentUserID(c *gin.Context) string {
	return c.GetString(middleware.UserIDKey)
}

// issueJWT creates a signed session token for the given user ID.
func (h *Handlers) issueJWT(userID string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": now.Unix(),
		"exp": now.Add(h.JWTExpiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.JWTSecret))
}

// getBudgets returns all budgets for a user, ordered by creation (saving first).
func getBudgets(ctx context.Context, q querier, userID string) ([]models.UserBudget, error) {
	rows, err := q.Query(ctx, `
		SELECT id, name, note, icon, color, percentage, amount
		FROM user_budgets
		WHERE user_id = $1
		ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	budgets := []models.UserBudget{}
	for rows.Next() {
		var b models.UserBudget
		if err := rows.Scan(&b.ID, &b.Name, &b.Note, &b.Icon, &b.Color, &b.Percentage, &b.Amount); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}
	return budgets, rows.Err()
}

// getBudget fetches a single budget owned by the user, or a notFound apiError.
func getBudget(ctx context.Context, q querier, userID string, id int64) (*models.UserBudget, error) {
	var b models.UserBudget
	err := q.QueryRow(ctx, `
		SELECT id, name, note, icon, color, percentage, amount
		FROM user_budgets
		WHERE id = $1 AND user_id = $2`, id, userID).
		Scan(&b.ID, &b.Name, &b.Note, &b.Icon, &b.Color, &b.Percentage, &b.Amount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, notFound("budget not found")
		}
		return nil, err
	}
	return &b, nil
}

// getSavingBudget fetches the user's protected "saving" bucket.
func getSavingBudget(ctx context.Context, q querier, userID string) (*models.UserBudget, error) {
	var b models.UserBudget
	err := q.QueryRow(ctx, `
		SELECT id, name, note, icon, color, percentage, amount
		FROM user_budgets
		WHERE user_id = $1 AND name = $2`, userID, savingBudgetName).
		Scan(&b.ID, &b.Name, &b.Note, &b.Icon, &b.Color, &b.Percentage, &b.Amount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("saving bucket missing for user %s", userID)
		}
		return nil, err
	}
	return &b, nil
}

// syncPercentages updates the percentages of the given non-saving budgets and
// recomputes the saving bucket's percentage so the total stays at 100%. It
// must be called inside the caller's transaction.
func syncPercentages(ctx context.Context, tx pgx.Tx, userID string, updated map[string]float64) error {
	sum := 0.0
	for idStr, pct := range updated {
		if pct < -percentageEpsilon {
			return badRequest("percentage for budget %s cannot be negative", idStr)
		}
		sum += pct
	}
	if sum > 100+percentageEpsilon {
		return badRequest("percentages must sum to at most 100%%")
	}

	for idStr, pct := range updated {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return badRequest("invalid budget id %q", idStr)
		}
		tag, err := tx.Exec(ctx, `
			UPDATE user_budgets SET percentage = $1
			WHERE id = $2 AND user_id = $3 AND name <> $4`, pct, id, userID, savingBudgetName)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return notFound("budget %d not found", id)
		}
	}

	savingPct := 100 - sum
	if savingPct < 0 {
		savingPct = 0
	}
	if _, err := tx.Exec(ctx, `
		UPDATE user_budgets SET percentage = $1
		WHERE user_id = $2 AND name = $3`, savingPct, userID, savingBudgetName); err != nil {
		return err
	}
	return nil
}

// insertBudget creates a new user_budgets row with amount=0 and returns its id.
func insertBudget(ctx context.Context, q querier, userID, name string, note *string, icon, color string, percentage float64) (int64, error) {
	var id int64
	err := q.QueryRow(ctx, `
		INSERT INTO user_budgets (user_id, name, note, icon, color, percentage, amount)
		VALUES ($1, $2, $3, $4, $5, $6, 0)
		RETURNING id`, userID, name, note, icon, color, percentage).Scan(&id)
	return id, err
}

// distributeAmounts proportionally splits amount across budgets by
// percentage, rounded to cents, assigning any rounding leftover to the
// saving bucket. budgets must include the saving bucket.
func distributeAmounts(budgets []models.UserBudget, amount float64) map[int64]float64 {
	allocations := make(map[int64]float64, len(budgets))
	var savingID int64
	var total float64
	for _, b := range budgets {
		alloc := math.Round(amount*b.Percentage/100*100) / 100
		allocations[b.ID] = alloc
		total += alloc
		if b.Name == savingBudgetName {
			savingID = b.ID
		}
	}

	leftover := math.Round((amount-total)*100) / 100
	if leftover != 0 {
		allocations[savingID] += leftover
	}
	return allocations
}
