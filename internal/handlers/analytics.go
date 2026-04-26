package handlers

import (
	"net/http"
	"strconv"
	"time"

	"expense-tracker/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnalyticsHandler struct {
	db *gorm.DB
}

func NewAnalyticsHandler(db *gorm.DB) *AnalyticsHandler {
	return &AnalyticsHandler{db: db}
}

type CategorySummary struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Total        float64 `json:"total"`
	Budget       float64 `json:"budget"`
	Overspent    bool    `json:"overspent"`
}

type MonthlySummary struct {
	Month      string            `json:"month"`
	Total      float64           `json:"total"`
	Count      int               `json:"count"`
	Categories []CategorySummary `json:"categories"`
}

func (h *AnalyticsHandler) GetMonthlySummary(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	if y := c.Query("year"); y != "" {
		parsed, err := strconv.Atoi(y)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
			return
		}
		year = parsed
	}
	if m := c.Query("month"); m != "" {
		parsed, err := strconv.Atoi(m)
		if err != nil || parsed < 1 || parsed > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
			return
		}
		month = parsed
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	var expenses []models.Expense
	h.db.Preload("Category").
		Where("user_id = ? AND date >= ? AND date < ?", userID, startDate, endDate).
		Find(&expenses)

	var total float64
	categoryMap := make(map[uint]*CategorySummary)

	for _, e := range expenses {
		total += e.Amount
		if _, ok := categoryMap[e.CategoryID]; !ok {
			categoryMap[e.CategoryID] = &CategorySummary{
				CategoryID:   e.CategoryID,
				CategoryName: e.Category.Name,
				Budget:       e.Category.Budget,
			}
		}
		categoryMap[e.CategoryID].Total += e.Amount
	}

	categories := make([]CategorySummary, 0, len(categoryMap))
	for _, cs := range categoryMap {
		cs.Overspent = cs.Budget > 0 && cs.Total > cs.Budget
		categories = append(categories, *cs)
	}

	c.JSON(http.StatusOK, MonthlySummary{
		Month:      startDate.Format("2006-01"),
		Total:      total,
		Count:      len(expenses),
		Categories: categories,
	})
}

func (h *AnalyticsHandler) GetCategoryBreakdown(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	type rawResult struct {
		CategoryID uint
		Total      float64
		Count      int
	}

	query := h.db.Model(&models.Expense{}).
		Select("category_id, SUM(amount) as total, COUNT(*) as count").
		Where("user_id = ?", userID).
		Group("category_id")

	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	var results []rawResult
	query.Scan(&results)

	type EnrichedResult struct {
		CategoryID   uint    `json:"category_id"`
		CategoryName string  `json:"category_name"`
		Total        float64 `json:"total"`
		Count        int     `json:"count"`
		Budget       float64 `json:"budget"`
		Overspent    bool    `json:"overspent"`
	}

	enriched := make([]EnrichedResult, 0, len(results))
	for _, r := range results {
		var cat models.Category
		h.db.First(&cat, r.CategoryID)
		enriched = append(enriched, EnrichedResult{
			CategoryID:   r.CategoryID,
			CategoryName: cat.Name,
			Total:        r.Total,
			Count:        r.Count,
			Budget:       cat.Budget,
			Overspent:    cat.Budget > 0 && r.Total > cat.Budget,
		})
	}

	c.JSON(http.StatusOK, enriched)
}
