package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kizoukun/codingtest/usecase"
	"github.com/kizoukun/codingtest/web"
)

func GetTodoController(w http.ResponseWriter, r *http.Request) {
	var req web.GetTodoRequest
	var response web.ResponseHttp
	ctx := r.Context()
	vars := mux.Vars(r)
	req.BoardID, _ = strconv.Atoi(vars["board_id"])

	authUsecase := usecase.NewTodoUsecase()
	authUsecase.GetTodoHandler(ctx, req, &response)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func AddTodoController(w http.ResponseWriter, r *http.Request) {
	var req web.TodoRequest
	var response web.ResponseHttp
	ctx := r.Context()
	vars := mux.Vars(r)
	req.BoardID, _ = strconv.Atoi(vars["board_id"])

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invalid request payload"
		json.NewEncoder(w).Encode(response)
		return
	}

	todoUsecase := usecase.NewTodoUsecase()
	todoUsecase.AddTodoHandler(ctx, req, &response)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func DeleteTodoController(w http.ResponseWriter, r *http.Request) {
	var req web.DeleteTodoRequest
	var response web.ResponseHttp
	vars := mux.Vars(r)
	req.ID = vars["id"]
	req.BoardID, _ = strconv.Atoi(vars["board_id"])
	ctx := r.Context()

	todoUsecase := usecase.NewTodoUsecase()
	todoUsecase.DeleteTodoHandler(ctx, req, &response)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func ToggleTodoController(w http.ResponseWriter, r *http.Request) {
	var req web.ToggleTodoRequest
	var response web.ResponseHttp
	ctx := r.Context()
	vars := mux.Vars(r)
	req.ID = vars["id"]
	req.BoardID, _ = strconv.Atoi(vars["board_id"])

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invalid request payload"
		json.NewEncoder(w).Encode(response)
		return
	}

	todoUsecase := usecase.NewTodoUsecase()
	todoUsecase.ToggleTodoHandler(ctx, req, &response)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}
