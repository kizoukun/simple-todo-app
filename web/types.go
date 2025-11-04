package web

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/kizoukun/codingtest/entity"
)

type TodoRequest struct {
	Task string `json:"task"`
}

type TodoGetResponse struct {
	Todos []entity.Todo `json:"todos"`
}

type ToggleTodoRequest struct {
	ID        string `json:"id"`
	Completed bool   `json:"completed"`
}

type DeleteTodoRequest struct {
	ID string `json:"id"`
}

type JwtClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}
