package controllers

import (
	"furniture-ecommerce/config"
	"furniture-ecommerce/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)
func AddToCart(c *gin.Context) {
    userID, exists := c.Get("userID") // Extract userID from JWT token
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var cart models.Cart
    if err := c.BindJSON(&cart); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    cart.UserID = userID.(uint) // Associate cart with the user

    if cart.ProductID == 0 || cart.Quantity == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "productid and quantity are required"})
        return
    }

    // Check the product exists
    var product models.Product
    if err := config.DB.Where("id = ?", cart.ProductID).First(&product).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Product does not exist"})
        return
    }

 
    var existingCart models.Cart
    if err := config.DB.Where("user_id = ? AND product_id = ?", cart.UserID, cart.ProductID).First(&existingCart).Error; err == nil {
        existingCart.Quantity += cart.Quantity
        config.DB.Save(&existingCart)
        c.JSON(http.StatusOK, gin.H{"message": "Cart updated successfully", "checkout": "/checkout"})
        return
    }

    
    if result := config.DB.Create(&cart); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add to cart"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":  "Added to cart successfully",
        "checkout": "/checkout",
    })
}


func RemoveFromCart(c *gin.Context) {
	userID, exists := c.Get("userID") // Extract userID from JWT token
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		log.Println(exists)
		return
	}

	cartID := c.Param("id")
	if result := config.DB.Where("id = ? AND user_id = ?", cartID, userID).Delete(&models.Cart{}); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not remove item from cart"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Removed from cart"})
}



func GetCart(c *gin.Context) {
	userID, exists := c.Get("userID") 
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var cart []models.Cart

	if result := config.DB.Where("user_id = ?", userID).Preload("Product").Find(&cart); result.Error != nil {
		log.Println("Error retrieving cart:", result.Error) 
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve cart items"})
		return
	}

	// Check if the cart is empty
	if len(cart) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Your cart is empty"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cart":    cart,
		"message": "Here are your cart items",
	})
}
