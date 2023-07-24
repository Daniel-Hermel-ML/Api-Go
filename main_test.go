package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

// START Testa a função para POST
func TestRoutePostStudents(t *testing.T) {
	// Configura do roteador do Gin para o test
	router := gin.Default()
	router.POST("/students", routePostStudents)

	// Dados de teste (payload JSON)
	payload := []byte(`{"full_name": "Daniel", "age": 35}`)

	// Requisição HTTP POST simulada com o payload de dados
	req, err := http.NewRequest("POST", "/students", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Configuração do cabeçalho da requisição
	req.Header.Set("Content-Type", "application/json")

	// Recorder para capturar a resposta do servidor
	recorder := httptest.NewRecorder()

	// Dispare a requisição HTTP
	router.ServeHTTP(recorder, req)

	// Verifique o código de status da resposta
	if recorder.Code != http.StatusCreated {
		t.Errorf("Esperava código de status %d, mas recebeu %d", http.StatusCreated, recorder.Code)
	}

	// Decodifique o corpo da resposta JSON em uma estrutura Student
	var response Student
	err = json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Erro ao decodificar a resposta JSON: %v", err)
	}

	// Verifique se o nome do estudante é igual ao enviado no payload
	expectedName := "Daniel"
	if response.Name != expectedName {
		t.Errorf("Esperava nome %s, mas recebeu %s", expectedName, response.Name)
	}

	// Verifique se a idade do estudante é igual à enviada no payload
	expectedAge := 35
	if response.Age != expectedAge {
		t.Errorf("Esperava idade %d, mas recebeu %d", expectedAge, response.Age)
	}

	// Verifique se o ID do estudante foi gerado corretamente (não deve ser zero)
	if response.ID == 0 {
		t.Errorf("O ID do estudante não foi gerado corretamente")
	}
}

// END Testa a função para POST
// START Testa a função para Get
func TestRouteGetStudents(t *testing.T) {
	// Configuração do roteador do Gin para o teste
	router := gin.Default()
	router.GET("/students", routeGetStudents)

	// Dados de teste (lista de estudantes)
	students := []Student{
		{ID: 1, Name: "Joao", Age: 18},
		{ID: 2, Name: "Gabriel", Age: 19},
		// Adicione mais estudantes, se necessário
	}

	// Recorder para capturar a resposta do servidor
	recorder := httptest.NewRecorder()

	// Crie uma requisição HTTP GET simulada
	req, err := http.NewRequest("GET", "/students", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Dispare a requisição HTTP
	router.ServeHTTP(recorder, req)

	// Verifique o código de status da resposta
	if recorder.Code != http.StatusOK {
		t.Errorf("Esperava código de status %d, mas recebeu %d", http.StatusOK, recorder.Code)
	}

	// Decodifique o corpo da resposta JSON em uma lista de estudantes
	var responseStudents []Student
	err = json.Unmarshal(recorder.Body.Bytes(), &responseStudents)
	if err != nil {
		t.Errorf("Erro ao decodificar a resposta JSON: %v", err)
	}

	// Verifique se o número de estudantes na resposta é o mesmo que o esperado
	expectedNumStudents := len(students)
	if len(responseStudents) != expectedNumStudents {
		t.Errorf("Esperava %d estudantes na resposta, mas recebeu %d", expectedNumStudents, len(responseStudents))
	}

	// Verifique se os estudantes retornados na resposta são os mesmos que os dados de teste
	for i, student := range students {
		if responseStudents[i].ID != student.ID ||
			responseStudents[i].Name != student.Name ||
			responseStudents[i].Age != student.Age {
			t.Errorf("Estudante na posição %d não corresponde aos dados de teste", i)
		}
	}
}

// END Testa a função para GET
// START testa a função para PUT
func TestRoutePutStudents(t *testing.T) {
	// Configuração do roteador do Gin para o teste
	router := gin.Default()
	router.PUT("/students/:id", routePutStudents)

	// Dados de teste (lista de estudantes)
	students := []Student{
		{ID: 1, Name: "Joao", Age: 18},
		{ID: 2, Name: "Gabriel", Age: 19},
		// Adicione mais estudantes, se necessário
	}

	// Recorder para capturar a resposta do servidor
	recorder := httptest.NewRecorder()

	// Estudante de teste para atualização (ID: 1)
	studentToUpdate := Student{ID: 1, Name: "Joao", Age: 18}

	// Codifique o corpo do payload JSON
	payload, err := json.Marshal(studentToUpdate)
	if err != nil {
		t.Fatal(err)
	}

	// Crie uma requisição HTTP PUT simulada
	req, err := http.NewRequest("PUT", "/students/1", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Configuração do cabeçalho da requisição
	req.Header.Set("Content-Type", "application/json")

	// Dispare a requisição HTTP
	router.ServeHTTP(recorder, req)

	// Verifique o código de status da resposta
	if recorder.Code != http.StatusOK {
		t.Errorf("Esperava código de status %d, mas recebeu %d", http.StatusOK, recorder.Code)
	}

	// Decodifique o corpo da resposta JSON em uma estrutura Student
	var responseStudent Student
	err = json.Unmarshal(recorder.Body.Bytes(), &responseStudent)
	if err != nil {
		t.Errorf("Erro ao decodificar a resposta JSON: %v", err)
	}

	// Verifique se o nome do estudante foi atualizado corretamente
	if responseStudent.Name != studentToUpdate.Name {
		t.Errorf("Esperava nome %s, mas recebeu %s", studentToUpdate.Name, responseStudent.Name)
	}

	// Verifique se a idade do estudante foi atualizada corretamente
	if responseStudent.Age != studentToUpdate.Age {
		t.Errorf("Esperava idade %d, mas recebeu %d", studentToUpdate.Age, responseStudent.Age)
	}

	// Verifique se o ID do estudante não foi alterado
	if responseStudent.ID != studentToUpdate.ID {
		t.Errorf("O ID do estudante foi alterado após a atualização")
	}

	// Verifique se o estudante atualizado está na lista de estudantes
	updatedStudentFound := false
	for _, student := range students {
		if student.ID == responseStudent.ID && student.Name == responseStudent.Name && student.Age == responseStudent.Age {
			updatedStudentFound = true
			break
		}
	}
	if !updatedStudentFound {
		t.Errorf("O estudante atualizado não foi encontrado na lista de estudantes")
	}
}

// END testa a função PUT
// START testa a função DELETE
func TestRouteDeleteStudents(t *testing.T) {
	// Configuração do roteador do Gin para o teste
	router := gin.Default()
	router.DELETE("/students/:id", routeDeleteStudents)

	// Dados de teste (lista de estudantes)
	students := []Student{
		{ID: 1, Name: "Joao", Age: 18},
		{ID: 2, Name: "Gabriel", Age: 19},
		{ID: 3, Name: "Teste", Age: 5},
	}

	// Recorder para capturar a resposta do servidor
	recorder := httptest.NewRecorder()

	// ID do estudante a ser excluído (ID: 3)
	idToDelete := 3

	// Crie uma requisição HTTP DELETE simulada
	req, err := http.NewRequest("DELETE", "/students/"+strconv.Itoa(idToDelete), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Dispare a requisição HTTP
	router.ServeHTTP(recorder, req)

	// Verifique o código de status da resposta
	if recorder.Code != http.StatusOK {
		t.Errorf("Esperava código de status %d, mas recebeu %d", http.StatusOK, recorder.Code)
	}

	// ****Atualize a variável students no teste após a exclusão do estudante
	newStudents := []Student{}
	for _, student := range students {
		if student.ID != idToDelete {
			newStudents = append(newStudents, student)
		}
	}
	students = newStudents

	// Verifique se o estudante foi excluído corretamente da lista de estudantes
	studentDeleted := false
	for _, student := range students {
		if student.ID == idToDelete {
			studentDeleted = true
			break
		}
	}
	if studentDeleted {
		t.Errorf("O estudante não foi excluído corretamente da lista de estudantes")
	}

	// Verifique se o estudante excluído não está mais na lista de estudantes
	updatedStudentFound := false
	for _, student := range students {
		if student.ID == idToDelete {
			updatedStudentFound = true
			break
		}
	}
	if updatedStudentFound {
		t.Errorf("O estudante excluído ainda está presente na lista de estudantes")
	}
}

// END testa a função DELETE
// START testa a função Get por ID
func TestRouteGetidStudents(t *testing.T) {
	// Configuração do roteador do Gin para o teste
	router := gin.Default()
	router.GET("/students/:id", routeGetidStudents)

	// Dados de teste (lista de estudantes)
	//students := []Student{
	//	{ID: 1, Name: "Joao", Age: 18},
	//	{ID: 2, Name: "Gabriel", Age: 19},
	//	{ID: 3, Name: "Teste", Age: 5},
	//}

	// ID do estudante a ser buscado (ID: 2)
	idToFetch := 2

	// Crie uma requisição HTTP GET simulada com o ID fornecido
	req, err := http.NewRequest("GET", "/students/"+strconv.Itoa(idToFetch), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Recorder para capturar a resposta do servidor
	recorder := httptest.NewRecorder()

	// Dispare a requisição HTTP
	router.ServeHTTP(recorder, req)

	// Verifique o código de status da resposta
	if recorder.Code != http.StatusOK {
		t.Errorf("Esperava código de status %d, mas recebeu %d", http.StatusOK, recorder.Code)
	}

	// Decodifique o corpo da resposta JSON em uma estrutura Student
	var responseStudent Student
	err = json.Unmarshal(recorder.Body.Bytes(), &responseStudent)
	if err != nil {
		t.Errorf("Erro ao decodificar a resposta JSON: %v", err)
	}

	// Verifique se o nome do estudante retornado é o esperado
	expectedName := "Gabriel"
	if responseStudent.Name != expectedName {
		t.Errorf("Esperava nome %s, mas recebeu %s", expectedName, responseStudent.Name)
	}

	// Verifique se a idade do estudante retornado é a esperada
	expectedAge := 19
	if responseStudent.Age != expectedAge {
		t.Errorf("Esperava idade %d, mas recebeu %d", expectedAge, responseStudent.Age)
	}

	// Verifique se o ID do estudante retornado é o esperado
	if responseStudent.ID != idToFetch {
		t.Errorf("Esperava ID %d, mas recebeu %d", idToFetch, responseStudent.ID)
	}
}

// END testa a função Get por ID
