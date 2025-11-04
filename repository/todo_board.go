package repository

import (
	"errors"
	"time"

	"github.com/kizoukun/codingtest/entity"
	"github.com/kizoukun/codingtest/mock"
)

type TodoBoardRepository struct {
	db *mock.MockDB[entity.TodoBoard]
}

func NewTodoBoardRepository() *TodoBoardRepository {
	dbBoard := mock.NewDb[entity.TodoBoard]("todo_boards.json")

	return &TodoBoardRepository{
		db: dbBoard,
	}
}

func (r *TodoBoardRepository) GetTodoBoard() ([]entity.TodoBoard, error) {
	data, err := r.db.GetData()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *TodoBoardRepository) CreateTodoBoard(todoBoard entity.TodoBoard) error {

	todoBoard.CreatedAt = time.Now()
	todoBoard.UpdatedAt = time.Now()
	err := r.db.InsertData(todoBoard)
	if err != nil {
		return err
	}
	return nil
}
func (r *TodoBoardRepository) UpdateTodoBoard(todoBoards []entity.TodoBoard) error {
	err := r.db.UpdateData(todoBoards)
	if err != nil {
		return err
	}
	return nil
}

func (r *TodoBoardRepository) GetById(boardId int) (entity.TodoBoard, error) {
	boards, err := r.GetTodoBoard()
	if err != nil {
		return entity.TodoBoard{}, err
	}
	for _, board := range boards {
		if board.ID == boardId {
			return board, nil
		}
	}
	return entity.TodoBoard{}, errors.New("not found")
}

func (r *TodoBoardRepository) GetByOwnerId(userID int) ([]entity.TodoBoard, error) {

	ownedBoards := make([]entity.TodoBoard, 0)
	// well in real case it should be in the db query to do this but this a mock up
	data, err := r.db.GetData()
	if err != nil {
		return ownedBoards, err
	}

	for _, board := range data {
		if board.CreatedBy == userID {
			ownedBoards = append(ownedBoards, board)
		}
	}

	return ownedBoards, nil
}
