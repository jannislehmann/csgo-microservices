package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID            uint64 `gorm:"primaryKey;autoIncrement:false"`
	ApiEnabled    bool
	ApiKey        string
	AuthCode      string
	LastShareCode string
}

func NewUser(id uint64) *User {
	return &User{
		ID: id,
	}
}
