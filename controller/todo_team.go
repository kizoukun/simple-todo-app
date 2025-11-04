package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kizoukun/codingtest/usecase"
	"github.com/kizoukun/codingtest/web"
)

func InviteTodoTeamController(w http.ResponseWriter, r *http.Request) {
	var req web.InviteTodoTeamsRequest
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

	todoTeamsUsecase := usecase.NewTodoTeamsUsecase()
	todoTeamsUsecase.InviteTodoTeamsHandler(ctx, req, &response)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}

func AcceptTodoTeamInviteController(w http.ResponseWriter, r *http.Request) {
	var req web.AcceptTodoTeamsInviteRequest
	var response web.ResponseHttp
	ctx := r.Context()

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invalid request payload"
		json.NewEncoder(w).Encode(response)
		return
	}

	todoTeamsUsecase := usecase.NewTodoTeamsUsecase()
	todoTeamsUsecase.AcceptTodoTeamsInvite(ctx, req, &response)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}
