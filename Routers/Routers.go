package Routers

import (
	"praktyka/Controllers"
	"praktyka/docs"
	"praktyka/Authentication"
	"net/http"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)


func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.POST("student", Authentication.TokenAuthMiddleware(),Controllers.AddNewStudent)
		v1.PUT("student/:id", Authentication.TokenAuthMiddleware(),Controllers.UpdateStudent)
		v1.DELETE("student/:id", Authentication.TokenAuthMiddleware(),Controllers.DeleteStudent)
		v1.GET("student", Controllers.GetStudents)
	}
	v2 := r.Group("/v2")
	{
		v2.POST("login", Controllers.UserLogin)
		v2.POST("logout", Authentication.TokenAuthMiddleware(), Controllers.UserLogout)
		v2.POST("register", Controllers.UserRegister)
		v2.POST("refresh", Controllers.TokenRefresh)
		v2.POST("role", Authentication.TokenAuthMiddleware(), Controllers.GetRole)
	}
	r.LoadHTMLGlob("sites/*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	docs.SwaggerInfo.BasePath = "/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
