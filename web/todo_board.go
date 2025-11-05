package web

import "github.com/kizoukun/codingtest/entity"

type CreateTodoBoardRequest struct {
	Title       string `json:"title" validate:"required,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

type CreateTodoBoardResponseData struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GetTodoBoardsRequest struct {
}

type GetTodoBoardsResponseData struct {
	OwnedBoards []entity.TodoBoard `json:"owned_boards"`
	TeamBoards  []entity.TodoBoard `json:"team_boards"`
}

type DeleteTodoBoardRequest struct {
	BoardID int `json:"board_id"`
}

type UpdateTodoBoardRequest struct {
	BoardID     int    `json:"board_id"`
	Title       string `json:"title" validate:"required,max=255"`
	Description string `json:"description" validate:"max=1000"`
}
