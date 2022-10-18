package Models

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Username string `json:"name"`
	HashedPassword string `json:"hashed_password"`
}
