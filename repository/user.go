package repository

import (
	"fmt"

	"github.com/kizoukun/codingtest/entity"
	"github.com/kizoukun/codingtest/mock"
)

var db = mock.NewDb[entity.User]("users.json")

func GetUsers() ([]entity.User, error) {
	data, err := db.GetData()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func CreateUser(user entity.User) error {

	err := db.InsertData(user)
	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) (*entity.User, error) {
	users, err := GetUsers()
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
