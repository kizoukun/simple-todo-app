package repository

import (
	"fmt"

	"github.com/kizoukun/codingtest/entity"
	"github.com/kizoukun/codingtest/mock"
)

type UserRepository struct {
	db *mock.MockDB[entity.User]
}

func NewUserRepository() *UserRepository {
	dbUser := mock.NewDb[entity.User]("users.json")
	return &UserRepository{
		db: dbUser,
	}
}

func (repo *UserRepository) GetUsers() ([]entity.User, error) {
	data, err := repo.db.GetData()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (repo *UserRepository) CreateUser(user entity.User) error {

	err := repo.db.InsertData(user)
	if err != nil {
		return err
	}
	return nil
}

func (repo *UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	users, err := repo.GetUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}
