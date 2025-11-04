package web

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/kizoukun/codingtest/entity"
)

type TodoRequest struct {
	BoardID int    `json:"board_id"`
	Task    string `json:"task"`
}

type TodoGetResponse struct {
	Todos []entity.Todo `json:"todos"`
}

type ToggleTodoRequest struct {
	BoardID   int    `json:"board_id"`
	ID        string `json:"id"`
	Completed bool   `json:"completed"`
}

type DeleteTodoRequest struct {
	BoardID int    `json:"board_id"`
	ID      string `json:"id"`
}

type JwtClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type GetTodoRequest struct {
	BoardID int `json:"board_id"`
}
