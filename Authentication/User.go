package Authentication

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"praktyka/ApiHelpers"
	"praktyka/Models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Ten cały plik powinien być w folderze "Controllers", nie ma co rozbijać kontrolerów na poszczególne foldery

func comparePasswords(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}

func Login(c *gin.Context) (interface{}, int) {
	var user Models.User
	err := c.BindJSON(&user)
	if err != nil {
		return user, 400
	}
	dbUser, err := Models.GetUserByUsername(user.Username)
	if err != nil {
		return "error occured with searching for user", 400
	}
	if !comparePasswords(user.Password, dbUser.Password) {
		return "incorrect passowrd", 400
	}
	ts, error := CreateToken(uint64(dbUser.ID), dbUser.Role)
	if error != nil {
		return error.Error(), 400
	}
	saveErr := Models.CreateAuth(dbUser.Role, ts)
	if saveErr != nil {
		return saveErr.Error(), 400
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	return tokens, 200
}

func CreateToken(id uint64, role string) (*Models.TokenDetails, error) {

	td := &Models.TokenDetails{}

	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = id
	atClaims["user_role"] = role
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	if err != nil {
		return nil, errors.New("access")
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = id
	rtClaims["user_role"] = role
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))

	if err != nil {
		return nil, errors.New("refresh")
	}

	return td, nil
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadata(r *http.Request) (*Models.AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		userRole := fmt.Sprintf("%.f", claims["user_role"])
		return &Models.AccessDetails{
			AccessUuid: accessUuid,
			UserRole:   userRole,
		}, nil
	}
	return nil, err
}

func Logout(c *gin.Context) (string, int) {
	au, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		return "unauthorized", 400
	}
	Models.DeleteAuth(au.AccessUuid)
	return "logout", 200
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := TokenValid(c.Request)
		if err != nil {
			ApiHelpers.RespondJSON(c, 400, "Token error")
			c.Abort()
			return
		}
		c.Next()
	}
}

func Refresh(c *gin.Context) (map[string]string, int) {
	refreshToken := c.Request.Header.Get("Refresh")

	if err := c.ShouldBindHeader(&refreshToken); err != nil {
		return nil, 400
	}
	fmt.Println(refreshToken)
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil

	})

	if err != nil {
		return nil, 400
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, 400
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			return nil, 400
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, 400
		}
		userRole := fmt.Sprintf("%.f", claims["user_role"])
		Models.DeleteAuth(refreshUuid)
		ts, createdError := CreateToken(userId, userRole)
		if createdError != nil {
			return nil, 400
		}
		saveError := Models.CreateAuth(userRole, ts)
		if saveError != nil {
			return nil, 400
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		return tokens, 200
	} else {
		return nil, 400
	}
}

func GetRole(c *gin.Context) string {
	tokenAuth, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		return ""
	}
	userRole := Models.FetchAuth(tokenAuth)

	if err != nil {
		return ""
	}
	return userRole
}

// Taka funkcja powinna być jako Middleware https://drstearns.github.io/tutorials/gomiddleware/
// W odpowiednim folderze "Middleware"
func RoleAuthentication(c *gin.Context, action string) (err error) {
	role := GetRole(c)
	switch role {
	case "administrator":
		return nil
	case "developer":
		if action == "update" {

			return nil
		}
		return errors.New("invalid role")
	}
	return errors.New("invalid role")
}
