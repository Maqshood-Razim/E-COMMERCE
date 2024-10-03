package models

import "gorm.io/gorm"


type Checkout struct {
	gorm.Model
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	UserID      uint   `json:"userid"`  // Foreign key to the User
	ProductID   uint   `json:"productid"` // Foreign key to the Product
	Price       uint   `json:"price"` 
	Product     Product `gorm:"foreignKey:ProductID"` 
	OrderDetails []OrderDetails `gorm:"foreignKey:CheckoutID" json:"order_details"`
}
