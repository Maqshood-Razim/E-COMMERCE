package models

type Payment struct {
	CheckoutID      uint    `json:"checkoutid"`        // this field for linking with checkout
	UserID          uint    `json:"userid"`           
	PaymentMode     string  `json:"payment_mode"`      
	Amount          float64 `json:"amount"`           
	Status          string  `json:"status" gorm:"default:'Pending'"`            
	StripePaymentID string  `json:"stripe_payment_id"` // Stripe Payment ID for tracking
	OrderID         uint    `json:"orderid"`           // Foreign key for linking to the OrderDetails table
}
