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
