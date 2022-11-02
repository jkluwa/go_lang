package Controllers

import (
	"praktyka/ApiHelpers"
	"praktyka/Authentication"
	"praktyka/Models"

	"github.com/gin-gonic/gin"
)

func UserLogin(c *gin.Context) {
	respond, status := Authentication.Login(c)
	ApiHelpers.RespondJSON(c, status, respond)

}

func UserRegister(c *gin.Context) {
	var user Models.User
	err := c.BindJSON(&user)
	if err != nil {
		ApiHelpers.RespondJSON(c, 400, err.Error())
		return
	}

	// Można uproscić sprawdź komentarz pod funkcją: Controllers/AddNewStudent
	err = Models.AddUser(&user)
	if err != nil {
		ApiHelpers.RespondJSON(c, 400, err.Error())
	} else {
		ApiHelpers.RespondJSON(c, 200, user)
	}
}

func UserLogout(c *gin.Context) {

	// Bez sensu jest odwołanie do osobnego folderu Authentication, pozbyć się go i tutaj wrzucić całą logikę funkcji do wylogowywania, to samo dla reszty na dole
	status, code := Authentication.Logout(c)
	ApiHelpers.RespondJSON(c, code, status)
}

func TokenRefresh(c *gin.Context) {
	status, code := Authentication.Refresh(c)
	ApiHelpers.RespondJSON(c, code, status)
}

func GetRole(c *gin.Context) {
	// Można uprosicić.
	role := Authentication.GetRole(c)
	if role == "" {
		ApiHelpers.RespondJSON(c, 400, "Error occured")
	} else {
		ApiHelpers.RespondJSON(c, 200, role)
	}
}
