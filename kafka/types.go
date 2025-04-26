package kafka


type CartEvent struct {
	UserID    uint   `json:"user_id"`
	ProductID uint   `json:"product_id"`
	Quantity  uint   `json:"quantity"`
	Price     uint   `json:"price"`
	Action    string `json:"action"` 
}

type OrderEvent struct {
	OrderID    uint    `json:"order_id"`
	UserID     uint    `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
}

type PaymentEvent struct {
	OrderID         uint    `json:"order_id"`
	UserID          uint    `json:"user_id"`
	Amount          float64 `json:"amount"`
	PaymentMethod   string  `json:"payment_method"`
	Status          string  `json:"status"`
	StripePaymentID string  `json:"stripe_payment_id,omitempty"`
}