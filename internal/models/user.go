package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string  `gorm:"size:35" json:"name"`
	Surname    string  `gorm:"size:35" json:"surname"`
	Username   string  `gorm:"uniqueIndex;size:35;not null" json:"username"`
	Email      string  `gorm:"uniqueIndex" json:"email"`
	Password   string  `gorm:"type:varchar(255)" json:"-"`
	TotalMoney float64 `json:"totalMoney"`
	Roles      []Role  `gorm:"many2many:user_roles"`

	// Define relationships
	FinanceHistory []FinanceRecord `gorm:"foreignKey:UserID"`
}
