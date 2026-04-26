package handlers

import (
	"net/http"
	"strconv"

	"expense-tracker/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	db *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		UserID: userID,
		Name:   req.Name,
		Budget: req.Budget,
	}

	if err := h.db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var categories []models.Category
	if err := h.db.Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var category models.Category
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Budget >= 0 {
		category.Budget = req.Budget
	}

	h.db.Save(&category)
	c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var count int64
	h.db.Model(&models.Expense{}).Where("category_id = ? AND user_id = ?", id, userID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete category with existing expenses"})
		return
	}

	result := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Category{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
