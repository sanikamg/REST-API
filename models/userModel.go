package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
}

type ID struct {
	gorm.Model
	Id uint `gorm:"unique"`
}
