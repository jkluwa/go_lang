package Controllers

import (
	"praktyka/ApiHelpers"
	"praktyka/Models"
	"github.com/gin-gonic/gin"
	"strconv"
)
// addStudent godoc
// @Summary      add student
// @Description  add
// @Tags         student
// @Accept       json
// @Produce      json
// @Param   	 Name      query   string     true  "Some name"
// @Param   	 Surname  query   string     true  "Some surname"
// @Param   	 Age       query   string     true  "Some age"
// @Success      200  {object}  string "ok"
// @Failure		 404  {object}	string "not ok"
// @Router       /student [post]
func AddNewStudent(c *gin.Context) {
	var student Models.Student
	c.BindJSON(&student)
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
// @Param   	 Name      query   string     true  "Some name"
// @Param   	 Surname   query   string     true  "Some surname"
// @Param   	 Age       query   string     true  "Some age"
// @Param   	 id      path   string     true  "Some id"
// @Success      200  {object}  string "ok"
// @Failure		 404  {object}	string "not ok"
// @Router       /student/{id} [put]
func UpdateStudent(c *gin.Context) {
	var student Models.Student
	c.BindJSON(&student)
	id, er := strconv.ParseUint(c.Params.ByName("id"), 10, 32)
	if er != nil {
		ApiHelpers.RespondJSON(c, 404, student)
	}
	student.ID = uint(id)
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
	id := c.Params.ByName("id")
	var student Models.Student = Models.GetStudent(id)
	err := Models.DeleteStudent(&student, id)
	if err != nil {
		ApiHelpers.RespondJSON(c, 404, &student)
	} else {
		ApiHelpers.RespondJSON(c, 200, &student)
	}
}
