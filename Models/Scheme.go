package Models

import (
	"time"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name string `json:"name"`
	Surname string `json:"surname"`
	Age string `json:"age"`
}

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Role string `json:"role"`
}

type TokenDetails struct {
	AccessToken		string
	RefreshToken	string
	AccessUuid		string
	RefreshUuid		string
	AtExpires		int64
	RtExpires		int64
}

type AccessDetails struct {
    AccessUuid string
	UserRole string
	exp time.Duration
}