package Routers

import (
	"praktyka/Controllers"
	"praktyka/docs"
	"net/http"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-contrib/cors"
)


func SetupRouter() *gin.Engine {
	r := gin.Default()
	
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOrigins = []string{"*"}
	// To be able to send tokens to the server.
	corsConfig.AllowCredentials = true

	// OPTIONS method for ReactJS
	corsConfig.AddAllowMethods("*")

	// Register the middleware
	r.Use(cors.New(corsConfig))

	v1 := r.Group("/v1")
	{
		v1.POST("student", Controllers.AddNewStudent)
		v1.PUT("student/:id", Controllers.UpdateStudent)
		v1.DELETE("student/:id", Controllers.DeleteStudent)
	}
	r.LoadHTMLGlob("sites/*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	docs.SwaggerInfo.BasePath = "/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
