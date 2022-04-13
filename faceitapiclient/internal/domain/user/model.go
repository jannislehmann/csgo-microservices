package user

import (
	"github.com/Cludch/csgo-microservices/shared/pkg/entity"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID         entity.ID `gorm:"primaryKey"`
	ApiEnabled bool
}

func NewUser(id entity.ID) *User {
	return &User{
		ID: id,
	}
}
