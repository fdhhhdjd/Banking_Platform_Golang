package user_services

import (
	"errors"

	"github.com/fdhhhdjd/Banking_Platform_Golang/internals/models"
)

func GetAllUsers() ([]models.User, error) {
	users := []models.User{
		{ID: 1, Name: "Nguyen Tien Tai", Email: "tai@example.com"},
	}

	if len(users) == 0 {
		return nil, errors.New("no users found")
	}

	return users, nil
}
