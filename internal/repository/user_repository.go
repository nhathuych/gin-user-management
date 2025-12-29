package repository

import (
	"gin-user-management/internal/model"
	"log"
)

type UserRepository struct {
	users []model.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make([]model.User, 0),
	}
}

func (ur *UserRepository) GetAll() {
	log.Println("GetAll called")
}

func (ur *UserRepository) Create() {
	log.Println("Create called")
}

func (ur *UserRepository) GetByUUID() {
	log.Println("GetByUUID called")
}

func (ur *UserRepository) Update() {
	log.Println("Update called")
}

func (ur *UserRepository) Delete() {
	log.Println("Delete called")
}
