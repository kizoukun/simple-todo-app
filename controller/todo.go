package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kizoukun/codingtest/usecase"
	"github.com/kizoukun/codingtest/web"
)

func GetTodoController(w http.ResponseWriter, r *http.Request) {
	var req interface{}
	var response web.ResponseHttp

	authUsecase := usecase.NewTodoUsecase()
	authUsecase.GetTodoHandler(req, &response)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func AddTodoController(w http.ResponseWriter, r *http.Request) {
	var req web.TodoRequest
	var response web.ResponseHttp

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invalid request payload"
		json.NewEncoder(w).Encode(response)
		return
	}

	todoUsecase := usecase.NewTodoUsecase()
	todoUsecase.AddTodoHandler(req, &response)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func DeleteTodoController(w http.ResponseWriter, r *http.Request) {
	var req web.DeleteTodoRequest
	var response web.ResponseHttp
	vars := mux.Vars(r)
	req.ID = vars["id"]

	todoUsecase := usecase.NewTodoUsecase()
	todoUsecase.DeleteTodoHandler(req, &response)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func ToggleTodoController(w http.ResponseWriter, r *http.Request) {
	var req web.ToggleTodoRequest
	var response web.ResponseHttp
	vars := mux.Vars(r)
	req.ID = vars["id"]

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invalid request payload"
		json.NewEncoder(w).Encode(response)
		return
	}

	todoUsecase := usecase.NewTodoUsecase()
	todoUsecase.ToggleTodoHandler(req, &response)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}
