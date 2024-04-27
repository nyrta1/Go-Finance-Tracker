package models

import "gorm.io/gorm"

//type RoleType string
//
//const (
//	USER  RoleType = "USER"
//	ADMIN RoleType = "ADMIN"
//)

type Role struct {
	gorm.Model
	Name  string `gorm:"size:35"`
	Users []User `gorm:"many2many:user_roles"`
}
