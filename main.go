package main

import (
	"praktyka/Config"
	"praktyka/Models"
	"praktyka/Routers"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var err error
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /v1

func main() {
	Config.DB, err = gorm.Open(sqlite.Open("school.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	Config.DB.AutoMigrate(&Models.Student{})

	r := Routers.SetupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
 