package Authentication

import (
	"praktyka/ApiHelpers"
	"fmt"
	"net/http"
	"os"
	"errors"
	"praktyka/Models"
	"strconv"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

var client *redis.Client



func RedisInit() {
	dsn := os.Getenv("REDIS_DSN")
	client = redis.NewClient(&redis.Options{
		Addr: dsn,
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func CreateAuth(id uint64, td *Models.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)

	now := time.Now()

	errAccess := client.Set(td.AccessUuid, id, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := client.Set(td.RefreshUuid, id, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func comparePasswords (password, hashedPassword string) (bool) {
	fmt.Println([]byte(password), " : ", []byte(hashedPassword))
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
	if(!comparePasswords(user.Password, dbUser.Password)) {
		return "incorrect passowrd", 400
	}
	ts, error := CreateToken(uint64(user.ID), user.Role)
	if error != nil {
		return error.Error(), 400
	}
	saveErr := CreateAuth(uint64(user.ID), ts)
	if saveErr != nil {
		return saveErr.Error(), 400
	}
	tokens := map[string]string{
		"access_token": ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	return tokens, 200
}

func CreateToken(id uint64, role string) (*Models.TokenDetails, error){

	td := &Models.TokenDetails{}

	td.AtExpires = time.Now().Add(time.Minute*15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour*24*7).Unix()
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

func VerifyToken(r *http.Request) (*jwt.Token, error){
	tokenString := r.Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func (token *jwt.Token) (interface{}, error) {
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
	if ok && token.Valid{
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		userRole := fmt.Sprintf("%.f", claims["user_role"])
		return &Models.AccessDetails{
			AccessUuid: accessUuid,
			UserRole: userRole,
			UserId: userId,
		}, nil
	}
	return nil, err
}

func fetchAuth(authD *Models.AccessDetails) (string, error) {
	userRole, err := client.Get(authD.UserRole).Result()
	if err != nil {
		return "", err
	}
	return userRole, nil
}

func Logout(c *gin.Context) (string, int) {
	au, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		return "unauthorized", 400
	}
	deleted, delErr := DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 {
		return "unauthorized", 400
	}
	return "logout", 200
}

func DeleteAuth(uuid string) (int64, error) {
	deleted, err := client.Del(uuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
	   err := TokenValid(c.Request)
	   if err != nil {
			fmt.Println(c.Request.Header.Get("Authorization"))
			ApiHelpers.RespondJSON(c, 400, "Token error")
		  c.Abort()
		  return 
	   }
	   c.Next()
	}
  }

func Refresh(c *gin.Context) (map[string]string, int){
	tokens := map[string]string{}
	if err := c.ShouldBindJSON(&tokens); err != nil {
		return nil, 400
	}
	refreshToken := tokens["refresh_token"]
	token, err := jwt.Parse(refreshToken, func (token *jwt.Token) (interface{}, error) {
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
		deleted, delError := DeleteAuth(refreshUuid)
		if delError != nil || deleted == 0 {
			return nil, 400
		}
		ts, createdError := CreateToken(userId, userRole)
		if createdError != nil {
			return nil, 400
		}
		saveError := CreateAuth(userId, ts)
		if saveError != nil {
			return nil, 400
		}
		tokens = map[string]string{
			"access_token": ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		return tokens, 200
	} else {
		return nil, 400
	}
}

func GetRole (c *gin.Context) string {
	tokenAuth, err := ExtractTokenMetadata(c.Request)
	if err != nil {
		return ""
	}
	userRole, err := fetchAuth(tokenAuth)

	if err != nil {
		return ""
	}
	return userRole
}

func RoleAuthentication (c *gin.Context, action string) (err error) {
	role := GetRole(c)
	switch role {
	case "administrator":
		return nil
	case "developer":
		if action == "update" {
			return nil
		}
		return err
	}
	return err
}