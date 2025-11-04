package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kizoukun/codingtest/usecase"
	"github.com/kizoukun/codingtest/web"
)

func GetTodoBoardController(w http.ResponseWriter, r *http.Request) {
	var req web.GetTodoBoardsRequest
	var response web.ResponseHttp
	ctx := r.Context()
	usecase := usecase.NewTodoBoardsUsecase()
	usecase.GetTodosByBoardHandler(ctx, req, &response)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func AddTodoBoardController(w http.ResponseWriter, r *http.Request) {
	var req web.CreateTodoBoardRequest
	var response web.ResponseHttp
	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invalid request payload"
		json.NewEncoder(w).Encode(response)
		return
	}

	todoBoardUsecase := usecase.NewTodoBoardsUsecase()
	todoBoardUsecase.CreateTodoBoardHandler(ctx, req, &response)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func DeleteTodoBoardController(w http.ResponseWriter, r *http.Request) {
	var req web.DeleteTodoBoardRequest
	var response web.ResponseHttp
	ctx := r.Context()
	vars := mux.Vars(r)
	req.BoardID, _ = strconv.Atoi(vars["id"])

	todoBoardUsecase := usecase.NewTodoBoardsUsecase()
	todoBoardUsecase.DeleteTodoBoardHandler(ctx, req, &response)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}
