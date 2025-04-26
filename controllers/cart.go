// controllers/cart.go
package controllers

import (
	"furniture-ecommerce/config"
	"furniture-ecommerce/kafka"
	"furniture-ecommerce/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartEvent struct {
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Quantity  uint     `json:"quantity"`
	Price     uint `json:"price"`
	Action    string  `json:"action"` 
}

func AddToCart(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var cart models.Cart
    if err := c.BindJSON(&cart); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    cart.UserID = userID.(uint)

    if cart.ProductID == 0 || cart.Quantity == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "productid and quantity are required"})
        return
    }

    var product models.Product
    if err := config.DB.Where("id = ?", cart.ProductID).First(&product).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Product does not exist"})
        return
    }

    var existingCart models.Cart
    if err := config.DB.Where("user_id = ? AND product_id = ?", cart.UserID, cart.ProductID).First(&existingCart).Error; err == nil {
        existingCart.Quantity += cart.Quantity
        config.DB.Save(&existingCart)
        
        // Publish Kafka event
        event := CartEvent{
            UserID:    cart.UserID,
            ProductID: cart.ProductID,
            Quantity: cart.Quantity,
            Price:     product.Price,
            Action:    "update",
        }
        if err := kafka.PublishEvent("cart_events", event); err != nil {
            log.Printf("Failed to publish cart event: %v", err)
        }
        
        c.JSON(http.StatusOK, gin.H{"message": "Cart updated successfully", "checkout": "/checkout"})
        return
    }

    if result := config.DB.Create(&cart); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add to cart"})
        return
    }

    // Publish Kafka event
    event := CartEvent{
        UserID:    cart.UserID,
        ProductID: cart.ProductID,
        Quantity: cart.Quantity,
        Price:     product.Price,
        Action:    "add",
    }
    if err := kafka.PublishEvent("cart_events", event); err != nil {
        log.Printf("Failed to publish cart event: %v", err)
    }

    c.JSON(http.StatusOK, gin.H{
        "message":  "Added to cart successfully",
        "checkout": "/checkout",
    })
}

func RemoveFromCart(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    cartID := c.Param("id")
    
    
    var cartItem models.Cart
    if err := config.DB.Where("id = ? AND user_id = ?", cartID, userID).First(&cartItem).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Item not found in cart"})
        return
    }
    
    var product models.Product
    if err := config.DB.First(&product, cartItem.ProductID).Error; err != nil {
        log.Printf("Product not found: %v", err)
    }

    if result := config.DB.Where("id = ? AND user_id = ?", cartID, userID).Delete(&models.Cart{}); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not remove item from cart"})
        return
    }

    event := CartEvent{
        UserID:    userID.(uint),
        ProductID: cartItem.ProductID,
        Quantity: cartItem.Quantity,
        Price:     product.Price,
        Action:    "remove",
    }
    if err := kafka.PublishEvent("cart_events", event); err != nil {
        log.Printf("Failed to publish cart event: %v", err)
    }

    c.JSON(http.StatusOK, gin.H{"message": "Removed from cart"})
}