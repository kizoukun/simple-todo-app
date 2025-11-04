package usecase

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/kizoukun/codingtest/entity"
	"github.com/kizoukun/codingtest/helper"
	"github.com/kizoukun/codingtest/mock"
	"github.com/kizoukun/codingtest/repository"
	"github.com/kizoukun/codingtest/web"
)

type TodoTeamsUsecase struct {
	teamsRepo  *repository.TodoTeamRepository
	boardsRepo *repository.TodoBoardRepository
	userRepo   *repository.UserRepository
	redis      *mock.MockRedis
}

func NewTodoTeamsUsecase() *TodoTeamsUsecase {
	return &TodoTeamsUsecase{
		teamsRepo:  repository.NewTodoTeamRepository(),
		boardsRepo: repository.NewTodoBoardRepository(),
		userRepo:   repository.NewUserRepository(),
		redis:      mock.NewMockRedis(),
	}
}

func (uc *TodoTeamsUsecase) AcceptTodoTeamsInvite(ctx context.Context, req web.AcceptTodoTeamsInviteRequest, response *web.ResponseHttp) {

	data := uc.redis.GetData(req.InviteCode)

	if data == nil {
		response.StatusCode = http.StatusNotFound
		response.Message = "Invalid or expired invite code"
		return
	}

	inviteData, ok := data.(web.BoardInviteData)
	if !ok {
		response.StatusCode = http.StatusNotFound
		response.Message = "Invalid or expired invite code"
		return
	}

	team, err := uc.teamsRepo.GetByUserId(inviteData.InvitedUserID)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to fetch user teams: " + err.Error()
		return
	}

	for _, t := range team {
		if t.BoardID == inviteData.BoardID {
			response.StatusCode = http.StatusBadRequest
			response.Message = "User already a member of the board"
			return
		}
	}

	newTeam := entity.TodoTeam{
		BoardID: inviteData.BoardID,
		UserID:  inviteData.InvitedUserID,
	}

	err = uc.teamsRepo.CreateTodoTeam(newTeam)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to accept invite: " + err.Error()
		return
	}

	uc.redis.DeleteData(req.InviteCode)

	response.StatusCode = http.StatusOK
	response.Message = "Invite accepted successfully"
	response.Success = true
}

func (uc *TodoTeamsUsecase) InviteTodoTeamsHandler(ctx context.Context, req web.InviteTodoTeamsRequest, response *web.ResponseHttp) {

	user, ok := helper.UserFromContext(ctx)
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

	if user.ID != board.CreatedBy {
		response.StatusCode = http.StatusForbidden
		response.Message = "Only board owner can invite members"
		return
	}

	if user.Email == req.Email {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Cannot invite yourself"
		return
	}

	invitedUser, err := uc.userRepo.GetUserByEmail(req.Email)

	if err != nil || invitedUser == nil {
		response.StatusCode = http.StatusBadRequest
		response.Message = "Invited user not found"
		return
	}

	team, err := uc.teamsRepo.GetByUserId(invitedUser.ID)
	if err != nil {
		response.StatusCode = http.StatusInternalServerError
		response.Message = "Failed to fetch user teams: " + err.Error()
		return
	}

	for _, t := range team {
		if t.BoardID == req.BoardID {
			response.StatusCode = http.StatusBadRequest
			response.Message = "User already a member of the board"
			return
		}
	}

	inviteData := web.BoardInviteData{
		BoardID:       req.BoardID,
		InvitedUserID: invitedUser.ID,
	}

	inviteCode := uuid.New().String()

	uc.redis.SetData(inviteCode, inviteData)

	response.StatusCode = http.StatusOK
	response.Data = map[string]string{
		"invite_code": inviteCode,
	}
	response.Message = "Invite created successfully"
	response.Success = true
}
