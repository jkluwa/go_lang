package Models

import (
	"errors"
	"praktyka/Config"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GetUser(id int) (user User) {
	// Brak error handlingu
	Config.DB.Where("id = ?", id).First(&user)
	return user
}

func GetUserByUsername(username string) (user User, err error) {
	// Tutaj jest bardzo ładny error handling
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
	// To samo co w funkcji FetchAuth
	Config.DB.Model(user).Where("username = ?", user.Username).Count(&count)

	if count > 0 {
		return errors.New("blad")
	}

	if err = Config.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func CreateAuth(role string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)

	now := time.Now()
	errAccess := Config.DB.Create(AccessDetails{td.AccessUuid, role, at.Sub(now)}).Error
	if errAccess != nil {
		return errAccess
	}

	errRefresh := Config.DB.Create(AccessDetails{td.RefreshUuid, role, rt.Sub(now)}).Error
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func FetchAuth(authD *AccessDetails) string {
	var count int64
	// GORM ma taką fajną funkcję jak .RowsAffected
	// Albo możesz po prostu zrobić strzał do bazy i szukać .Error
	// Do poprawki, nie potrzeba takiej dziwnej logiki, dodatkowo brakuje samego error handlingu
	Config.DB.Model(authD).Where("access_uuid = ?", authD.AccessUuid).Count(&count)
	if count != 0 {
		Config.DB.Model(authD).Where("access_uuid = ?", authD.AccessUuid).First(&authD)

		return authD.UserRole
	}
	return ""
}

func DeleteAuth(uuid string) {
	// Brak error handlingu
	Config.DB.Where("access_uuid = ?", uuid).Delete(AccessDetails{})
}
