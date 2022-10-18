package Controllers

import (
	"praktyka/ApiHelpers"
	"praktyka/Models"
	"crypto/md5"
	"encoding/hex"
	"github.com/gin-gonic/gin"
)
// addStudent godoc
// @Summary      add student
// @Description  add
// @Tags         student
// @Accept       json
// @Produce      json
// @Param   	 Username      query   string     true  "Some string"
// @Success      200  {object}  string "ok"
// @Failure		 404  {object}	string "not ok"
// @Router       /student [post]
func AddNewStudent(c *gin.Context) {
	var student Models.Student
	c.BindJSON(&student)
	student.Username = c.Query("Username")
	hash := md5.Sum([]byte(c.Query("Password")))
	student.HashedPassword =  hex.EncodeToString(hash[:])
	err := Models.AddNewStudent(&student)
	if err != nil {
		ApiHelpers.RespondJSON(c, 404, student)
	} else {
		ApiHelpers.RespondJSON(c, 200, student)
	}
}
// updateStudent godoc
// @Summary      update student
// @Description  update
// @Tags         student
// @Accept       json
// @Produce      json
// @Param   	 Username      query   string     true  "Some string"
// @Param   	 id      path   string     true  "Some id"
// @Success      200  {object}  string "ok"
// @Failure		 404  {object}	string "not ok"
// @Router       /student/{id} [put]
func UpdateStudent(c *gin.Context) {
	var student Models.Student
	student.Username = c.DefaultQuery("Username", Models.GetStudent(c.Params.ByName("id")).Username)
	if c.DefaultQuery("Password", "") != "" {
		hash := md5.Sum([]byte(c.DefaultQuery("Password", "")))	
		student.HashedPassword = hex.EncodeToString(hash[:])
	} else {
		student.HashedPassword = Models.GetStudent(c.Params.ByName("id")).HashedPassword
	}
	err := Models.UpdateStudent(&student, c.Params.ByName("id"))
	if err != nil {
		ApiHelpers.RespondJSON(c, 404, student)
	} else {
		ApiHelpers.RespondJSON(c, 200, student)
	}
}
// deleteStudent godoc
// @Summary      delete student
// @Description  delete
// @Tags         student
// @Accept       json
// @Produce      json
// @Success      200  {object}  string "ok"
// @Failure		 404  {object}	string "not ok"
// @Param   	 id      path   string     true  "Some id"
// @Router       /student/{id} [delete]
func DeleteStudent(c *gin.Context) {
	var student Models.Student
	id := c.Params.ByName("id")
	err := Models.DeleteStudent(&student, id)
	if err != nil {
		ApiHelpers.RespondJSON(c, 404, student)
	} else {
		ApiHelpers.RespondJSON(c, 200, student)
	}
}
