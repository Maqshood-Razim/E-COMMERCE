package models

type User struct {
	ID       uint    `json:"id" gorm:"primary_key"`
	Username string `gorm:"unique" json:"username"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}
