package controllers

import (
	"furniture-ecommerce/config"
	"furniture-ecommerce/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BrowseProducts(c *gin.Context) {
	var products []models.Product
	if result := config.DB.Find(&products); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

func SearchProduct(c *gin.Context) {
	var input struct {
		Name string `json:"name"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var products []models.Product
	
	log.Printf("Searching for product with name: %s", input.Name)

	result := config.DB.Where("name LIKE ?", "%"+input.Name+"%").Find(&products)

	if result.Error != nil {
		
		log.Printf("Database query error: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error", "details": result.Error.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find products"})
		return
	}

	c.JSON(http.StatusOK, products)
}