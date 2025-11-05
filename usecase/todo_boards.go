package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/kizoukun/codingtest/entity"
	"github.com/kizoukun/codingtest/helper"
	"github.com/kizoukun/codingtest/repository"
	"github.com/kizoukun/codingtest/web"
)

type TodoBoardsUsecase struct {
	boardsRepo *repository.TodoBoardRepository
	teamRepo   *repository.TodoTeamRepository
}

func NewTodoBoardsUsecase() *TodoBoardsUsecase {
	return &TodoBoardsUsecase{
		boardsRepo: repository.NewTodoBoardRepository(),
		teamRepo:   repository.NewTodoTeamRepository(),
	}
}

func (uc *TodoBoardsUsecase) GetTodosByBoardHandler(context context.Context, req web.GetTodoBoardsRequest, response *web.ResponseHttp) {

	user, ok := helper.UserFromContext(context)
	if !ok || user == nil {
		response.StatusCode = http.StatusUnauthorized
		response.Message = "Unauthorized"
		return
	}

	owned, err := uc.boardsRepo.GetByOwnerId(user.ID)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to fetch todo boards: " + err.Error()
		return
	}

	teams, err := uc.teamRepo.GetByUserId(user.ID)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to fetch todo teams: " + err.Error()
		return
	}

	teamBoards := make([]entity.TodoBoard, 0)

	for _, team := range teams {
		board, err := uc.boardsRepo.GetById(team.BoardID)
		if err != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Message = "Failed to fetch todo board for team: " + err.Error()
			return
		}
		teamBoards = append(teamBoards, board)
	}

	response.StatusCode = http.StatusOK
	response.Data = web.GetTodoBoardsResponseData{
		OwnedBoards: owned,
		TeamBoards:  teamBoards,
	}
	response.Message = "Todo boards fetched successfully"
	response.Success = true

}

func (uc *TodoBoardsUsecase) CreateTodoBoardHandler(context context.Context, req web.CreateTodoBoardRequest, response *web.ResponseHttp) {
	user, ok := helper.UserFromContext(context)
	if !ok || user == nil {
		response.StatusCode = http.StatusUnauthorized
		response.Message = "Unauthorized"
		return
	}

	board := entity.TodoBoard{
		Title:       req.Title,
		Description: req.Description,
		CreatedBy:   user.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := uc.boardsRepo.CreateTodoBoard(board)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to create todo board: " + err.Error()
		return
	}

	response.StatusCode = http.StatusCreated
	response.Data = web.CreateTodoBoardResponseData{
		ID:          board.ID,
		Title:       board.Title,
		Description: board.Description,
	}
	response.Message = "Todo board created successfully"
	response.Success = true
}

func (uc *TodoBoardsUsecase) DeleteTodoBoardHandler(context context.Context, req web.DeleteTodoBoardRequest, response *web.ResponseHttp) {
	user, ok := helper.UserFromContext(context)
	if !ok || user == nil {
		response.StatusCode = http.StatusUnauthorized
		response.Message = "Unauthorized"
		return
	}

	board, err := uc.boardsRepo.GetById(req.BoardID)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invalid Board ID: " + err.Error()
		return
	}

	if board.CreatedBy != user.ID {
		response.StatusCode = http.StatusForbidden
		response.Message = "You do not have permission to delete this board"
		return
	}

	allBoards, err := uc.boardsRepo.GetTodoBoard()
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to fetch todo boards: " + err.Error()
		return
	}

	updatedBoards := make([]entity.TodoBoard, 0)
	for _, b := range allBoards {
		if b.ID != req.BoardID {
			updatedBoards = append(updatedBoards, b)
		}
	}

	err = uc.boardsRepo.UpdateTodoBoard(updatedBoards)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to delete todo board: " + err.Error()
		return
	}

	response.StatusCode = http.StatusOK
	response.Message = "Todo board deleted successfully"
	response.Success = true
}

func (uc *TodoBoardsUsecase) UpdateTodoBoardHandler(context context.Context, req web.UpdateTodoBoardRequest, response *web.ResponseHttp) {
	user, ok := helper.UserFromContext(context)
	if !ok || user == nil {
		response.StatusCode = http.StatusUnauthorized
		response.Message = "Unauthorized"
		return
	}

	board, err := uc.boardsRepo.GetById(req.BoardID)
	if err != nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invalid Board ID: " + err.Error()
		return
	}

	if board.CreatedBy != user.ID {
		response.StatusCode = http.StatusForbidden
		response.Message = "You do not have permission to update this board"
		return
	}

	board.Title = req.Title
	board.Description = req.Description
	board.UpdatedAt = time.Now()

	boards, err := uc.boardsRepo.GetTodoBoard()
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to fetch todo boards: " + err.Error()
		return
	}

	// update the board in the list
	for i, b := range boards {
		if b.ID == board.ID {
			boards[i] = board
			break
		}
	}

	err = uc.boardsRepo.UpdateTodoBoard(boards)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to update todo board: " + err.Error()
		return
	}

	response.StatusCode = http.StatusOK
	response.Data = board
	response.Message = "Todo board updated successfully"
	response.Success = true
}
