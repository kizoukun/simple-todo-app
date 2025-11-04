package repository

import (
	"github.com/kizoukun/codingtest/entity"
	"github.com/kizoukun/codingtest/mock"
)

type TodoTeamRepository struct {
	db *mock.MockDB[entity.TodoTeam]
}

func NewTodoTeamRepository() *TodoTeamRepository {
	dbTeam := mock.NewDb[entity.TodoTeam]("todo_teams.json")

	return &TodoTeamRepository{
		db: dbTeam,
	}
}

func (r *TodoTeamRepository) GetTodoTeam() ([]entity.TodoTeam, error) {
	data, err := r.db.GetData()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *TodoTeamRepository) CreateTodoTeam(todoTeam entity.TodoTeam) error {

	err := r.db.InsertData(todoTeam)
	if err != nil {
		return err
	}
	return nil
}

func (r *TodoTeamRepository) UpdateTodoTeam(todoTeams []entity.TodoTeam) error {
	err := r.db.UpdateData(todoTeams)
	if err != nil {
		return err
	}
	return nil
}

func (r *TodoTeamRepository) GetByUserId(userId int) ([]entity.TodoTeam, error) {
	userTeams := make([]entity.TodoTeam, 0)
	teams, err := r.GetTodoTeam()
	if err != nil {
		return nil, err
	}
	for _, team := range teams {
		if team.UserID == userId {
			userTeams = append(userTeams, team)
		}
	}
	return userTeams, nil
}
