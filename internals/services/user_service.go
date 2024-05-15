package services

import "github.com/fdhhhdjd/Banking_Platform_Golang/internals/models"

func GetAllUsers() []models.User {
	return []models.User{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
	}
}
