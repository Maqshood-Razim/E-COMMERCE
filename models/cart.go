package models

type Cart struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	UserID    uint    `json:"userid"`
	ProductID uint    `json:"productid"`
	Quantity  uint    `json:"quantity"`
	Product   Product `gorm:"foreignKey:ProductID;references:ID"`
}
