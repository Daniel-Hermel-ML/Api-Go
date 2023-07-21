package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"full_name"`
	Age  int    `json:"age"`
}

var Students = []Student{
	Student{ID: 1, Name: "Joao", Age: 18},
	Student{ID: 2, Name: "Gabriel", Age: 19},
}

func routerHearth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Conectado...",
	})
	c.Done()
}

// START GET
func routeGetStudents(c *gin.Context) {
	c.JSON(http.StatusOK, Students)
}

// END GET
// START POST
func routePostStudents(c *gin.Context) {
	var student Student

	err := c.Bind(&student)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possível obter o payload - 37",
		})
		return
	}

	student.ID = Students[len(Students)-1].ID + 1
	Students = append(Students, student)

	c.JSON(http.StatusCreated, student)
}

// END POST
// START PUT
func routePutStudents(c *gin.Context) {
	var studentPayload Student
	var studentLocal Student
	var newStudents []Student

	err := c.BindJSON(&studentPayload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possível obter o Payload - 56",
		})
		return
	}

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possível obter o Payload - 64",
		})
		return
	}

	for _, studentElement := range Students {
		if studentElement.ID == id {
			studentLocal = studentElement
		}

	}
	if studentLocal.ID == 0 {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message_error": "Não foi possível encontrar o estudante - 78",
			})
			return
		}
	}
	studentLocal.Name = studentPayload.Name
	studentLocal.Age = studentPayload.Age

	for _, studentElement := range Students {
		if id == studentElement.ID {
			newStudents = append(newStudents, studentLocal)
		} else {
			newStudents = append(newStudents, studentElement)
		}
	}
	Students = newStudents

	c.JSON(http.StatusOK, studentLocal)

}

// END PUT
// START Delete
func routeDeleteStudents(c *gin.Context) {
	var newStudents []Student

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possível obter o id. 111",
		})
		return
	}

	for _, studentElement := range Students {
		if studentElement.ID != id {
			newStudents = append(newStudents, studentElement)
		}
	}

	Students = newStudents

	c.JSON(http.StatusOK, gin.H{
		"message": "Estudante excluído com sucesso!",
	})
}

// END Delete
// START Get por id
func routeGetidStudents(c *gin.Context) {
	var student Student

	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possível obter o id.",
		})
		return
	}

	for _, studentElement := range Students {
		if studentElement.ID == id {
			student = studentElement
		}
	}
	if student.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message_error": "Não foi possível encontrar o estudante",
		})
		return

	}

	c.JSON(http.StatusOK, student)
}

// END Get por id
// START func main
func main() {
	service := gin.Default()

	getRoutes(service)

	service.Run()
}

// END func main
// START rotas
func getRoutes(c *gin.Engine) *gin.Engine {
	c.GET("/heart", routerHearth)

	groupStudents := c.Group("/students")
	groupStudents.GET("/", routeGetStudents)
	groupStudents.POST("/", routePostStudents)
	groupStudents.PUT("/:id", routePutStudents)
	groupStudents.DELETE("/:id", routeDeleteStudents)
	groupStudents.GET("/:id", routeGetidStudents)

	return c
}

//END rotas
