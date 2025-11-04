package web

type AcceptTodoTeamsInviteRequest struct {
	InviteCode string `json:"invite_code" validate:"required,uuid4"`
}

type BoardInviteData struct {
	BoardID       int `json:"board_id"`
	InvitedUserID int `json:"user_id"`
}

type InviteTodoTeamsRequest struct {
	BoardID int    `json:"board_id" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}
