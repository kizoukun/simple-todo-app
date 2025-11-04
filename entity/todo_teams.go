package entity

type TodoTeam struct {
	ID      int `json:"id"`
	BoardID int `json:"board_id"`
	UserID  int `json:"user_id"`
}
