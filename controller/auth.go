package controller

import (
	"encoding/json"
	"net/http"

	"github.com/kizoukun/codingtest/usecase"
	"github.com/kizoukun/codingtest/web"
)

func LoginController(w http.ResponseWriter, r *http.Request) {
	var req web.LoginRequest
	var response web.ResponseHttp

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invalid request payload"
		json.NewEncoder(w).Encode(response)
		return
	}

	authUsecase := usecase.NewAuthUsecase()
	authUsecase.AuthLoginHandler(req, &response)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)

}

func RegisterController(w http.ResponseWriter, r *http.Request) {
	var req web.RegisterRequest
	var response web.ResponseHttp

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invalid request payload"
		json.NewEncoder(w).Encode(response)
		return
	}

	authUsecase := usecase.NewAuthUsecase()
	authUsecase.AuthRegisterHandler(req, &response)

	w.WriteHeader(response.StatusCode)
	json.NewEncoder(w).Encode(response)
}
