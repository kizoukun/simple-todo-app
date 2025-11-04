package usecase

import (
	"net/http"
	"strconv"
	"time"

	"github.com/kizoukun/codingtest/entity"
	"github.com/kizoukun/codingtest/repository"
	"github.com/kizoukun/codingtest/web"
)

type TodoUsecase struct {
	todoRepo *repository.TodoRepository
}

func NewTodoUsecase() *TodoUsecase {
	return &TodoUsecase{
		todoRepo: repository.NewTodoRepository(),
	}
}

func (uc *TodoUsecase) GetTodoHandler(req any, response *web.ResponseHttp) {

	data, err := uc.todoRepo.GetTodo()
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to fetch todos" + err.Error()
		return
	}

	response.StatusCode = http.StatusOK
	response.Data = data
	response.Message = "Todos fetched successfully"
	response.Success = true
}

func (uc *TodoUsecase) AddTodoHandler(req web.TodoRequest, response *web.ResponseHttp) {
	err := uc.todoRepo.CreateTodo(entity.Todo{
		Task: req.Task,
	})
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to create todo: " + err.Error()
		return
	}

	response.StatusCode = http.StatusCreated
	response.Data = req
	response.Message = "Todo created successfully"
	response.Success = true
}

func (uc *TodoUsecase) DeleteTodoHandler(req web.DeleteTodoRequest, response *web.ResponseHttp) {

	todos, err := uc.todoRepo.GetTodo()
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to fetch todos: " + err.Error()
		return
	}
	for index, todo := range todos {
		if strconv.Itoa(todo.ID) == req.ID {
			todos = append(todos[:index], todos[index+1:]...)
			response.Success = true
			response.StatusCode = http.StatusNoContent
			response.Message = "Success delete todo"
			uc.todoRepo.UpdateTodo(todos)
			return
		}
	}

	response.StatusCode = http.StatusNotFound
	response.Message = "Todo Not Found"

}

func (uc *TodoUsecase) ToggleTodoHandler(req web.ToggleTodoRequest, response *web.ResponseHttp) {

	todos, err := uc.todoRepo.GetTodo()
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to fetch todos: " + err.Error()
		return
	}
	for index, todo := range todos {
		if strconv.Itoa(todo.ID) == req.ID {
			todos[index].Completed = req.Completed
			todos[index].UpdatedAt = time.Now()
			response.Success = true
			response.StatusCode = http.StatusOK
			response.Message = "Success toggle todo"
			response.Data = todos[index]
			uc.todoRepo.UpdateTodo(todos)
			return
		}
	}

	response.StatusCode = http.StatusNotFound
	response.Message = "Todo Not Found"
}
