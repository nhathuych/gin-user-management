package service

import (
	"gin-user-management/internal/repository"
)

type userService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetAll() {
	us.repo.GetAll()
}

func (us *userService) Create() {
	us.repo.Create()
}

func (us *userService) GetByUUID() {
	us.repo.GetByUUID()
}

func (us *userService) Update() {
	us.repo.Update()
}

func (us *userService) Delete() {
	us.repo.Delete()
}
