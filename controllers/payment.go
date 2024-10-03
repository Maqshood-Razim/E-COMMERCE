package controllers

import (
	"furniture-ecommerce/config"
	"furniture-ecommerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

func ProcessPayment(c *gin.Context) {
	
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	
	var payment models.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Ensure the user provides a checkoutID
	if payment.CheckoutID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Checkout ID is required"})
		return
	}
	if payment.PaymentMode != "stripe" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid payment method",
			"message": "Stripe is the only payment method",
		})
		return
	}


	var order models.OrderDetails
	if err := config.DB.Where("checkout_id = ?", payment.CheckoutID).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Checkout not found"})
		return
	}

	
	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order is not pending"})
		return
	}

	// Check if the payment amount matches the order's total price
	if payment.Amount < order.TotalPrice {
		remaining := order.TotalPrice - payment.Amount
		c.JSON(http.StatusBadRequest, gin.H{
			"error":       "Insufficient payment",
			"remaining":   remaining,
			"total_price": order.TotalPrice,
		})
		return
	}

	// Initialize Stripe with the secret key
	stripe.Key = "sk_test_51Q3VaJP8UEcGEhHsoMSaJV7OzABvuJDya5Vy0CHib5Z6V0Dee4MVVBAzx95NwOpMXQ5z0ylaCwC03xtVy2lteyfc00pGTq7MwR"

	// Create a PaymentIntent on Stripe
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(payment.Amount * 100)), 
		Currency: stripe.String("inr"),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card", 
		}),
	}

	// Create the payment intent
	pi, err := paymentintent.New(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment intent"})
		return
	}

	// Immediately confirm the payment intent on the backend
	_, err = paymentintent.Confirm(
		pi.ID,
		&stripe.PaymentIntentConfirmParams{
			PaymentMethod: stripe.String("pm_card_visa"), 
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm payment intent"})
		return
	}

	// Update the payment record as successful
	payment.UserID = userID.(uint)
	payment.Status = "succeeded"
	payment.StripePaymentID = pi.ID

	// Save the payment in the database
	if result := config.DB.Create(&payment); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not process payment"})
		return
	}

	
	order.Status = "paid"
	if result := config.DB.Save(&order); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update order status"})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"message":       "Payment successful",
		"checkout_id":   payment.CheckoutID,
		"order_details": "/order/details",
	})
}