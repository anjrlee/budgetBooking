package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"budgetbooking/internal/models"
)

// trendPoint mirrors the frontend's MonthlyTrendPoint type.
type trendPoint struct {
	Month string  `json:"month"`
	Total float64 `json:"total"`
}

// spendingSummaryRow mirrors the frontend's SpendingSummaryRow type.
type spendingSummaryRow struct {
	BudgetID       *int64  `json:"budget_id"`
	BudgetName     string  `json:"budget_name"`
	Icon           *string `json:"icon,omitempty"`
	Color          *string `json:"color,omitempty"`
	TotalSpent     float64 `json:"total_spent"`
	CurrentBalance float64 `json:"current_balance"`
}

// parseMonthRange parses a "YYYY-MM" query param into a [start, end) range.
func parseMonthRange(month string) (time.Time, time.Time, error) {
	start, err := time.Parse("2006-01", month)
	if err != nil {
		return time.Time{}, time.Time{}, badRequest("invalid month %q, expected YYYY-MM", month)
	}
	return start, start.AddDate(0, 1, 0), nil
}

// parseYearRange parses a "YYYY" query param into a [start, end) range.
func parseYearRange(yearStr string) (time.Time, time.Time, error) {
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return time.Time{}, time.Time{}, badRequest("invalid year %q, expected YYYY", yearStr)
	}
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	return start, start.AddDate(1, 0, 0), nil
}

// queryBookings returns budget_bookings rows for the user within [start, end),
// optionally filtered to a single budget and/or limited.
func queryBookings(ctx context.Context, q querier, userID string, isIncome bool, start, end time.Time, budgetID *int64, limit *int) ([]models.BudgetBooking, error) {
	sql := `
		SELECT id, budget_id, budget_name, is_income, amount, note, created_at
		FROM budget_bookings
		WHERE user_id = $1 AND is_income = $2 AND created_at >= $3 AND created_at < $4`
	args := []interface{}{userID, isIncome, start, end}

	if budgetID != nil {
		args = append(args, *budgetID)
		sql += fmt.Sprintf(" AND budget_id = $%d", len(args))
	}
	sql += " ORDER BY created_at DESC"
	if limit != nil {
		args = append(args, *limit)
		sql += fmt.Sprintf(" LIMIT $%d", len(args))
	}

	rows, err := q.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := []models.BudgetBooking{}
	for rows.Next() {
		var b models.BudgetBooking
		if err := rows.Scan(&b.ID, &b.BudgetID, &b.BudgetName, &b.IsIncome, &b.Amount, &b.Note, &b.CreatedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, rows.Err()
}

// querySpendingSummary aggregates spending within [start, end) grouped by budget.
func querySpendingSummary(ctx context.Context, q querier, userID string, start, end time.Time) ([]spendingSummaryRow, error) {
	rows, err := q.Query(ctx, `
		SELECT bb.budget_id, COALESCE(bb.budget_name, ub.name), ub.icon, ub.color,
		       SUM(bb.amount), COALESCE(ub.amount, 0)
		FROM budget_bookings bb
		LEFT JOIN user_budgets ub ON ub.id = bb.budget_id
		WHERE bb.user_id = $1 AND bb.is_income = false AND bb.created_at >= $2 AND bb.created_at < $3
		GROUP BY bb.budget_id, COALESCE(bb.budget_name, ub.name), ub.icon, ub.color, ub.amount
		ORDER BY SUM(bb.amount) DESC`, userID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	summary := []spendingSummaryRow{}
	for rows.Next() {
		var row spendingSummaryRow
		if err := rows.Scan(&row.BudgetID, &row.BudgetName, &row.Icon, &row.Color, &row.TotalSpent, &row.CurrentBalance); err != nil {
			return nil, err
		}
		summary = append(summary, row)
	}
	return summary, rows.Err()
}

// dailyTrend returns one point per day of [start, end) with that day's total spending for budgetID.
func dailyTrend(ctx context.Context, q querier, userID string, budgetID int64, start, end time.Time) ([]trendPoint, error) {
	rows, err := q.Query(ctx, `
		SELECT created_at::date, SUM(amount)
		FROM budget_bookings
		WHERE user_id = $1 AND budget_id = $2 AND is_income = false AND created_at >= $3 AND created_at < $4
		GROUP BY created_at::date`, userID, budgetID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totals := map[string]float64{}
	for rows.Next() {
		var day time.Time
		var total float64
		if err := rows.Scan(&day, &total); err != nil {
			return nil, err
		}
		totals[day.Format("2006-01-02")] = total
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	points := []trendPoint{}
	for cur := start; cur.Before(end); cur = cur.AddDate(0, 0, 1) {
		key := cur.Format("2006-01-02")
		points = append(points, trendPoint{Month: key, Total: totals[key]})
	}
	return points, nil
}

// monthlyTrend returns one point per month of the given year with that month's total spending for budgetID.
func monthlyTrend(ctx context.Context, q querier, userID string, budgetID int64, year int) ([]trendPoint, error) {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	rows, err := q.Query(ctx, `
		SELECT date_trunc('month', created_at), SUM(amount)
		FROM budget_bookings
		WHERE user_id = $1 AND budget_id = $2 AND is_income = false AND created_at >= $3 AND created_at < $4
		GROUP BY date_trunc('month', created_at)`, userID, budgetID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totals := map[string]float64{}
	for rows.Next() {
		var month time.Time
		var total float64
		if err := rows.Scan(&month, &total); err != nil {
			return nil, err
		}
		totals[month.Format("2006-01")] = total
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	points := []trendPoint{}
	for m := 0; m < 12; m++ {
		key := start.AddDate(0, m, 0).Format("2006-01")
		points = append(points, trendPoint{Month: key, Total: totals[key]})
	}
	return points, nil
}

// availablePeriods returns the distinct months (YYYY-MM) and years with matching booking records.
func availablePeriods(ctx context.Context, q querier, userID string, isIncome bool, budgetID *int64) ([]string, []int, error) {
	args := []interface{}{userID, isIncome}
	filter := ""
	if budgetID != nil {
		args = append(args, *budgetID)
		filter = fmt.Sprintf(" AND budget_id = $%d", len(args))
	}

	monthRows, err := q.Query(ctx, `
		SELECT DISTINCT to_char(created_at, 'YYYY-MM') AS period
		FROM budget_bookings
		WHERE user_id = $1 AND is_income = $2`+filter+`
		ORDER BY period`, args...)
	if err != nil {
		return nil, nil, err
	}
	months := []string{}
	for monthRows.Next() {
		var m string
		if err := monthRows.Scan(&m); err != nil {
			monthRows.Close()
			return nil, nil, err
		}
		months = append(months, m)
	}
	monthRows.Close()
	if err := monthRows.Err(); err != nil {
		return nil, nil, err
	}

	yearRows, err := q.Query(ctx, `
		SELECT DISTINCT EXTRACT(YEAR FROM created_at)::int AS period
		FROM budget_bookings
		WHERE user_id = $1 AND is_income = $2`+filter+`
		ORDER BY period`, args...)
	if err != nil {
		return nil, nil, err
	}
	years := []int{}
	for yearRows.Next() {
		var y int
		if err := yearRows.Scan(&y); err != nil {
			yearRows.Close()
			return nil, nil, err
		}
		years = append(years, y)
	}
	yearRows.Close()
	if err := yearRows.Err(); err != nil {
		return nil, nil, err
	}

	return months, years, nil
}

// budgetIDQuery parses the required "budget_id" query parameter.
func budgetIDQuery(c *gin.Context) (int64, error) {
	idStr := c.Query("budget_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, badRequest("invalid budget_id %q", idStr)
	}
	return id, nil
}

// ShowSpendingReport handles GET /report/showSpendingReport.
func (h *Handlers) ShowSpendingReport(c *gin.Context) {
	userID := currentUserID(c)
	start, end, err := parseMonthRange(c.Query("month"))
	if err != nil {
		respondError(c, err)
		return
	}

	var limit *int
	if limitStr := c.Query("limit"); limitStr != "" {
		n, err := strconv.Atoi(limitStr)
		if err != nil || n < 0 {
			respondError(c, badRequest("invalid limit %q", limitStr))
			return
		}
		limit = &n
	}

	ctx := c.Request.Context()
	summary, err := querySpendingSummary(ctx, h.DB, userID, start, end)
	if err != nil {
		respondError(c, err)
		return
	}
	transactions, err := queryBookings(ctx, h.DB, userID, false, start, end, nil, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"summary": summary, "transactions": transactions})
}

// ShowBudgetReport handles GET /report/showBudgetReport.
func (h *Handlers) ShowBudgetReport(c *gin.Context) {
	userID := currentUserID(c)
	budgetID, err := budgetIDQuery(c)
	if err != nil {
		respondError(c, err)
		return
	}
	start, end, err := parseMonthRange(c.Query("month"))
	if err != nil {
		respondError(c, err)
		return
	}

	ctx := c.Request.Context()
	trend, err := dailyTrend(ctx, h.DB, userID, budgetID, start, end)
	if err != nil {
		respondError(c, err)
		return
	}
	transactions, err := queryBookings(ctx, h.DB, userID, false, start, end, &budgetID, nil)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"trend": trend, "transactions": transactions})
}

// ShowIncomeReport handles GET /report/showIncomeReport.
func (h *Handlers) ShowIncomeReport(c *gin.Context) {
	userID := currentUserID(c)
	start, end, err := parseMonthRange(c.Query("month"))
	if err != nil {
		respondError(c, err)
		return
	}

	var limit *int
	if limitStr := c.Query("limit"); limitStr != "" {
		n, err := strconv.Atoi(limitStr)
		if err != nil || n < 0 {
			respondError(c, badRequest("invalid limit %q", limitStr))
			return
		}
		limit = &n
	}

	rows, err := queryBookings(c.Request.Context(), h.DB, userID, true, start, end, nil, limit)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"rows": rows})
}

// AvailablePeriods handles GET /report/availablePeriods.
func (h *Handlers) AvailablePeriods(c *gin.Context) {
	userID := currentUserID(c)

	var isIncome bool
	switch c.Query("type") {
	case "spending":
		isIncome = false
	case "income":
		isIncome = true
	default:
		respondError(c, badRequest("invalid type %q, expected spending or income", c.Query("type")))
		return
	}

	var budgetID *int64
	if idStr := c.Query("budget_id"); idStr != "" {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			respondError(c, badRequest("invalid budget_id %q", idStr))
			return
		}
		budgetID = &id
	}

	months, years, err := availablePeriods(c.Request.Context(), h.DB, userID, isIncome, budgetID)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"months": months, "years": years})
}

// ShowAnnualSpendingReport handles GET /report/showAnnualSpendingReport.
func (h *Handlers) ShowAnnualSpendingReport(c *gin.Context) {
	userID := currentUserID(c)
	start, end, err := parseYearRange(c.Query("year"))
	if err != nil {
		respondError(c, err)
		return
	}

	ctx := c.Request.Context()
	summary, err := querySpendingSummary(ctx, h.DB, userID, start, end)
	if err != nil {
		respondError(c, err)
		return
	}
	transactions, err := queryBookings(ctx, h.DB, userID, false, start, end, nil, nil)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"summary": summary, "transactions": transactions})
}

// ShowAnnualBudgetReport handles GET /report/showAnnualBudgetReport.
func (h *Handlers) ShowAnnualBudgetReport(c *gin.Context) {
	userID := currentUserID(c)
	budgetID, err := budgetIDQuery(c)
	if err != nil {
		respondError(c, err)
		return
	}
	start, end, err := parseYearRange(c.Query("year"))
	if err != nil {
		respondError(c, err)
		return
	}

	ctx := c.Request.Context()
	trend, err := monthlyTrend(ctx, h.DB, userID, budgetID, start.Year())
	if err != nil {
		respondError(c, err)
		return
	}
	transactions, err := queryBookings(ctx, h.DB, userID, false, start, end, &budgetID, nil)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"trend": trend, "transactions": transactions})
}

// ShowAnnualIncomeReport handles GET /report/showAnnualIncomeReport.
func (h *Handlers) ShowAnnualIncomeReport(c *gin.Context) {
	userID := currentUserID(c)
	start, end, err := parseYearRange(c.Query("year"))
	if err != nil {
		respondError(c, err)
		return
	}

	rows, err := queryBookings(c.Request.Context(), h.DB, userID, true, start, end, nil, nil)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"rows": rows})
}
