package controllers

import (
	"furniture-ecommerce/config"
	"furniture-ecommerce/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Checkout(c *gin.Context) {
	// Retrieve user ID and other details from the JSON request
	var checkoutRequest struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
		UserID      uint   `json:"userid"`
		ProductID   uint   `json:"productid"`
		Price       uint   `json:"price"`
	}

	// Bind the JSON request to the struct
	if err := c.ShouldBindJSON(&checkoutRequest); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}


	var user models.User
	if result := config.DB.Where("id = ?", checkoutRequest.UserID).First(&user); result.Error != nil {
		log.Printf("User not found: %v", result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
		return
	}

	// Retrieve the user's cart items
	var cart []models.Cart
	if result := config.DB.Where("user_id = ?", checkoutRequest.UserID).Find(&cart); result.Error != nil || len(cart) == 0 {
		log.Printf("Cart retrieval error or cart empty: %v", result.Error)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	
	log.Printf("Processing checkout for user ID: %d with %d items in cart\n", checkoutRequest.UserID, len(cart))

	// Calculate the total price of the products in the cart
	var totalPrice uint
	var correctPrice uint
	for _, cartItem := range cart {
		var product models.Product
		if err := config.DB.First(&product, cartItem.ProductID).Error; err != nil {
			log.Printf("Product not found in cart: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found in cart"})
			return
		}
		totalPrice += product.Price * cartItem.Quantity

		// Check if this is the product the user is checking out
		if product.ID == checkoutRequest.ProductID {
			correctPrice = product.Price * cartItem.Quantity
		}
	}

	
	log.Printf("Total price calculated: %d\n", totalPrice)

	// Ensure the correct product price is calculated
	if correctPrice == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Product not found in cart",
			"message": "Please add the product to your cart",
		})
		return
	}

	// Compare the user-provided price with the correct price of the product
	if checkoutRequest.Price != correctPrice {
		remainingAmount := int(correctPrice) - int(checkoutRequest.Price)
		if remainingAmount < 0 {
			remainingAmount = 0 //avoid negative values
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":           "Incorrect price",
			"remaining_price": remainingAmount,
			"correct_price":   correctPrice,
			"message":         "Please enter the correct price for the product",
		})
		return
	}

	// Create a new Checkout object
	checkout := models.Checkout{
		Name:        checkoutRequest.Name,
		PhoneNumber: checkoutRequest.PhoneNumber,
		Address:     checkoutRequest.Address,
		UserID:      checkoutRequest.UserID,
		ProductID:   checkoutRequest.ProductID,
		Price:       totalPrice, 
	}

	// Save checkout details to the database
	if result := config.DB.Create(&checkout); result.Error != nil {
		log.Printf("Error saving checkout to database: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save checkout details"})
		return
	}

	orderDetails := models.OrderDetails{
		CheckoutID: checkout.ID,
		Quantity:   1, 
		Price:      float64(totalPrice),
		UserId:     checkoutRequest.UserID,
		TotalPrice: float64(totalPrice),
		Status:     "pending",
		Address:    checkoutRequest.Address,
		Products:   []models.Product{{ID: checkoutRequest.ProductID}},
	}

	if result := config.DB.Create(&orderDetails); result.Error != nil {
		log.Printf("Error saving order details to database: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save order details"})
		return
	}


	log.Printf("Checkout saved successfully: %v", checkout)

	
	c.JSON(http.StatusOK, gin.H{
		"message":      "Checkout successful",
		"total_price":  totalPrice,
		"make_payment": "/payment", 
		"checkout_id":  checkout.ID,
	})
}
