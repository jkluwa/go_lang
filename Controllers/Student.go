package Controllers

import (
	"praktyka/ApiHelpers"
	"praktyka/Authentication"
	"praktyka/Models"
	"strconv"

	"github.com/gin-gonic/gin"
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
// @Failure		 400  {object}	string "not ok"
// @Router       /student [post]
func AddNewStudent(c *gin.Context) {
	// Nie widzę w funkcji Role Authentication czegoś takiego jak "add"
	// Sprawdzanie czy ktoś ma dostęp do danego endpointa powinno być w Middlewarze, nie tutaj w każdym z controllerów
	// Jeżeli ktoś nie ma dostępu, to middleware zamyka request i zwraca odpowiedni błąd sprawdź komentarz pod funkcją: Authenticaiton/RoleAuthentication
	err := Authentication.RoleAuthentication(c, "add")
	if err != nil {
		ApiHelpers.RespondJSON(c, 400, "invalid role")
		return
	}
	var student Models.Student
	c.BindJSON(&student)

	// Nie potrzebna komplikacja

	// err = Models.AddNewStudent(&student)
	// if err != nil {
	// 	ApiHelpers.RespondJSON(c, 400, student)
	// } else {
	// 	ApiHelpers.RespondJSON(c, 200, student)
	// }

	// Nie lepiej tak?
	if err := Models.AddNewStudent(&student); err != nil {
		ApiHelpers.RespondJSON(c, 400, student)
		return
	}
	ApiHelpers.RespondJSON(c, 200, student)
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
// @Failure		 400  {object}	string "not ok"
// @Router       /student/{id} [put]
func UpdateStudent(c *gin.Context) {
	// Sprawdzanie czy ktoś ma dostęp do danego endpointa powinno być w Middlewarze, nie tutaj w każdym z controllerów
	// Jeżeli ktoś nie ma dostępu, to middleware zamyka request i zwraca odpowiedni błąd sprawdź komentarz pod funkcją: Authenticaiton/RoleAuthentication
	err := Authentication.RoleAuthentication(c, "update")
	if err != nil {
		ApiHelpers.RespondJSON(c, 400, "invalid role")
		return
	}
	var student Models.Student
	c.BindJSON(&student)
	id, err := strconv.ParseUint(c.Params.ByName("id"), 10, 32)
	if err != nil {
		ApiHelpers.RespondJSON(c, 400, student)
	}
	// Po co podwójnie podajesz ID usera? Najpierw do structa, a potem do samego requestu do bazy
	student.ID = uint(id)
	err = Models.UpdateStudent(&student, c.Params.ByName("id"))
	if err != nil {
		ApiHelpers.RespondJSON(c, 400, student)
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
// @Failure		 400  {object}	string "not ok"
// @Param   	 id      path   string     true  "Some id"
// @Router       /student/{id} [delete]
func DeleteStudent(c *gin.Context) {
	// Sprawdzanie czy ktoś ma dostęp do danego endpointa powinno być w Middlewarze, nie tutaj w każdym z controllerów
	// Jeżeli ktoś nie ma dostępu, to middleware zamyka request i zwraca odpowiedni błąd sprawdź komentarz pod funkcją: Authenticaiton/RoleAuthentication
	err := Authentication.RoleAuthentication(c, "delete")
	if err != nil {
		ApiHelpers.RespondJSON(c, 400, "invalid role")
		return
	}
	id := c.Params.ByName("id")
	var student Models.Student = Models.GetStudent(id)

	// Można uprościć, sprawdź funkcję AddNewStudent
	err = Models.DeleteStudent(&student, id)
	if err != nil {
		ApiHelpers.RespondJSON(c, 400, &student)
	} else {
		ApiHelpers.RespondJSON(c, 200, &student)
	}
}

// BRAKUJE KOMENTARZY SWAGGERA
func GetStudents(c *gin.Context) {

	// Można uprosić, sprawdź funkcję AddNewStudent
	students, err := Models.GetStudents()
	if err != nil {
		ApiHelpers.RespondJSON(c, 400, "Error occured")
	} else {
		ApiHelpers.RespondJSON(c, 200, students)
	}
}
