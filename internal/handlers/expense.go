package handlers

import (
	"net/http"
	"strconv"
	"time"

	"expense-tracker/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExpenseHandler struct {
	db *gorm.DB
}

func NewExpenseHandler(db *gorm.DB) *ExpenseHandler {
	return &ExpenseHandler{db: db}
}

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req models.CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	if err := h.db.Where("id = ? AND user_id = ?", req.CategoryID, userID).First(&category).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
		return
	}

	expense := models.Expense{
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        date,
	}

	if err := h.db.Create(&expense).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}

	h.db.Preload("Category").First(&expense, expense.ID)
	c.JSON(http.StatusCreated, expense)
}

func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	query := h.db.Where("user_id = ?", userID).Preload("Category")

	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("date >= ?", startDate)
	}
	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("date <= ?", endDate)
	}

	var expenses []models.Expense
	if err := query.Order("date desc").Find(&expenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}

	c.JSON(http.StatusOK, expenses)
}

func (h *ExpenseHandler) GetExpense(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var expense models.Expense
	if err := h.db.Preload("Category").Where("id = ? AND user_id = ?", id, userID).First(&expense).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	c.JSON(http.StatusOK, expense)
}

func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var expense models.Expense
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&expense).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	var req models.UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.CategoryID != 0 {
		var category models.Category
		if err := h.db.Where("id = ? AND user_id = ?", req.CategoryID, userID).First(&category).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
			return
		}
		expense.CategoryID = req.CategoryID
	}
	if req.Amount != 0 {
		expense.Amount = req.Amount
	}
	if req.Description != "" {
		expense.Description = req.Description
	}
	if req.Date != "" {
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD"})
			return
		}
		expense.Date = date
	}

	h.db.Save(&expense)
	h.db.Preload("Category").First(&expense, expense.ID)
	c.JSON(http.StatusOK, expense)
}

func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Expense{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted"})
}
