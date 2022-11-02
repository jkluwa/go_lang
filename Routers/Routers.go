package Routers

import (
	"net/http"
	"praktyka/Authentication"
	"praktyka/Controllers"
	"praktyka/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.POST("student", Authentication.TokenAuthMiddleware(), Controllers.AddNewStudent)
		v1.PUT("student/:id", Authentication.TokenAuthMiddleware(), Controllers.UpdateStudent)
		v1.DELETE("student/:id", Authentication.TokenAuthMiddleware(), Controllers.DeleteStudent)
		// Jeżeli bierzemy listę z bazy, to nazwyamy endpoint "GET /students"
		// Poprosę dorobić endpoint na GET dla pojedyńczego usera "GET /student/student:id"
		// https://blog.dreamfactory.com/best-practices-for-naming-rest-api-endpoints/
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
