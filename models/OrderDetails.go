package models

type OrderDetails struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CheckoutID uint      `json:"checkoutid"`    
	Quantity   uint      `json:"quantity"`
	Price      float64   `json:"price"`
	UserId     uint      `json:"userid"`
	TotalPrice float64   `json:"totalprice"`
	Status     string    `json:"status"`
	Address    string    `json:"address"`
	PaymentID  uint      `json:"paymentid"`
	Products   []Product `gorm:"many2many:order_products;" json:"products"` 
	Checkout   Checkout  `gorm:"foreignKey:CheckoutID" json:"checkout"`
	
}
