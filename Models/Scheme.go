package Models

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name string `json:"name"`
	Surname string `json:"surname"`
	Age string `json:"age"`
}
