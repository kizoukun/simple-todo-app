package usecase

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kizoukun/codingtest/web"
)

var todos = []web.Todo{}
var nextID = 1

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	response := web.TodoGetResponse{Todos: todos}
	json.NewEncoder(w).Encode(response)
}

func AddTodoHandler(w http.ResponseWriter, r *http.Request) {
	var newTodo web.Todo
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		errorResponse := web.ErrorResponse{
			Error: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	newTodo.ID = nextID
	nextID++
	todos = append(todos, newTodo)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for index, todo := range todos {
		if strconv.Itoa(todo.ID) == id {
			todos = append(todos[:index], todos[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	errorResponse := web.ErrorResponse{
		Error: "Todo Not Found",
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(errorResponse)
}

func ToggleTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var todoRequest web.ToggleTodoRequest
	err := json.NewDecoder(r.Body).Decode(&todoRequest)
	if err != nil {
		errorResponse := web.ErrorResponse{
			Error: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	for index, todo := range todos {
		if strconv.Itoa(todo.ID) == id {
			todos[index].Completed = todoRequest.Completed
			json.NewEncoder(w).Encode(todos[index])
			return
		}
	}
	errorResponse := web.ErrorResponse{
		Error: "Todo Not Found",
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(errorResponse)
}
