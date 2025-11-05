package usecase

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/kizoukun/codingtest/entity"
	"github.com/kizoukun/codingtest/helper"
	"github.com/kizoukun/codingtest/repository"
	"github.com/kizoukun/codingtest/web"
)

type TodoUsecase struct {
	todoRepo      *repository.TodoRepository
	todoBoardRepo *repository.TodoBoardRepository
}

func NewTodoUsecase() *TodoUsecase {
	return &TodoUsecase{
		todoRepo:      repository.NewTodoRepository(),
		todoBoardRepo: repository.NewTodoBoardRepository(),
	}
}

func (uc *TodoUsecase) GetTodoHandler(context context.Context, req web.GetTodoRequest, response *web.ResponseHttp) {

	user, ok := helper.UserFromContext(context)
	if !ok || user == nil {
		response.StatusCode = http.StatusUnauthorized
		response.Message = "Unauthorized"
		return
	}

	board, err := uc.todoBoardRepo.GetById(req.BoardID)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invalid Board ID: " + err.Error()
		return
	}

	data, err := uc.todoRepo.GetByBoardId(req.BoardID)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to fetch todos" + err.Error()
		return
	}

	response.StatusCode = http.StatusOK
	response.Data = map[string]interface{}{
		"board": board,
		"todos": data,
	}
	response.Message = "Todos fetched successfully"
	response.Success = true
}

func (uc *TodoUsecase) AddTodoHandler(context context.Context, req web.TodoRequest, response *web.ResponseHttp) {

	user, ok := helper.UserFromContext(context)
	if !ok || user == nil {
		response.StatusCode = http.StatusUnauthorized
		response.Message = "Unauthorized"
		return
	}

	if _, err := uc.todoBoardRepo.GetById(req.BoardID); err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invalid Board ID: " + err.Error()
		return
	}

	err := uc.todoRepo.CreateTodo(entity.Todo{
		Task:      req.Task,
		CreatedBy: user.ID,
		BoardId:   req.BoardID,
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

func (uc *TodoUsecase) DeleteTodoHandler(context context.Context, req web.DeleteTodoRequest, response *web.ResponseHttp) {

	if _, err := uc.todoBoardRepo.GetById(req.BoardID); err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invalid Board ID: " + err.Error()
		return
	}

	todos, err := uc.todoRepo.GetByBoardId(req.BoardID)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to fetch todos: " + err.Error()
		return
	}

	for index, todo := range todos {
		if strconv.Itoa(todo.ID) == req.ID {
			todos = append(todos[:index], todos[index+1:]...)
			response.Success = true
			response.StatusCode = http.StatusOK
			response.Message = "Success delete todo"
			uc.todoRepo.UpdateTodo(todos)
			return
		}
	}

	response.StatusCode = http.StatusNotFound
	response.Message = "Todo Not Found"

}

func (uc *TodoUsecase) ToggleTodoHandler(context context.Context, req web.ToggleTodoRequest, response *web.ResponseHttp) {

	if _, err := uc.todoBoardRepo.GetById(req.BoardID); err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invalid Board ID: " + err.Error()
		return
	}

	todos, err := uc.todoRepo.GetByBoardId(req.BoardID)
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
