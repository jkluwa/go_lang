package Routers

import (
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
		v1.POST("student", Controllers.AddNewStudent)
		v1.PATCH("student/:id", Controllers.UpdateStudent)
		v1.DELETE("student/:id", Controllers.DeleteStudent)
	}

	docs.SwaggerInfo.BasePath = "/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
