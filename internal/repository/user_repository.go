package repository

import (
	"gin-user-management/internal/model"
	"log"
)

type userRepository struct {
	users []model.User
}

func NewUserRepository() UserRepository {
	return &userRepository{
		users: make([]model.User, 0),
	}
}

func (ur *userRepository) GetAll() {
	log.Println("GetAll called")
}

func (ur *userRepository) Create() {
	log.Println("Create called")
}

func (ur *userRepository) GetByUUID() {
	log.Println("GetByUUID called")
}

func (ur *userRepository) Update() {
	log.Println("Update called")
}

func (ur *userRepository) Delete() {
	log.Println("Delete called")
}
