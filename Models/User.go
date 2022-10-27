package Models

import (
	//"encoding/hex"
	"errors"
	"praktyka/Config"

	"golang.org/x/crypto/bcrypt"
)

func GetUser(id int) (user User) {
	Config.DB.Where("id = ?", id).First(&user)
	return user
}

func GetUserByUsername(username string) (user User, err error) {
	if ok := Config.DB.Where("username = ?", username).First(&user).Error; ok != nil {
		return user, err
	}
	return user, nil
}

func AddUser(user *User) (err error) {
	
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error with hashing")
	}
	user.Password = string(hash)
	var count int64
	Config.DB.Model(user).Where("username = ?", user.Username).Count(&count)
	
	if count > 0 {
		return errors.New("blad")
	}
	

	if err = Config.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}