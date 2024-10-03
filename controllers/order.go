package controllers

import (
	"furniture-ecommerce/config"
	"furniture-ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OrderDetails(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	
	var order models.OrderDetails
	if err := config.DB.Preload("Products").Where("user_id = ?", userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}


	
	c.JSON(http.StatusOK, gin.H{
		"message": "Order details retrieved successfully",
		"order":   order,
	})
}
