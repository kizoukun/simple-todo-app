package repository

import (
	"time"

	"github.com/kizoukun/codingtest/entity"
	"github.com/kizoukun/codingtest/mock"
)

type TodoRepository struct {
	db *mock.MockDB[entity.Todo]
}

func NewTodoRepository() *TodoRepository {
	dbTodo := mock.NewDb[entity.Todo]("todos.json")
	return &TodoRepository{
		db: dbTodo,
	}
}

func (r *TodoRepository) GetTodo() ([]entity.Todo, error) {
	data, err := r.db.GetData()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *TodoRepository) CreateTodo(todo entity.Todo) error {

	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	err := r.db.InsertData(todo)
	if err != nil {
		return err
	}
	return nil
}

func (r *TodoRepository) UpdateTodo(todos []entity.Todo) error {
	err := r.db.UpdateData(todos)
	if err != nil {
		return err
	}
	return nil
}
