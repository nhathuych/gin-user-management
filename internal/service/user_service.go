package service

import (
	"gin-user-management/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) GetAll() {
	us.repo.GetAll()
}

func (us *UserService) Create() {
	us.repo.Create()
}

func (us *UserService) GetByUUID() {
	us.repo.GetByUUID()
}

func (us *UserService) Update() {
	us.repo.Update()
}

func (us *UserService) Delete() {
	us.repo.Delete()
}
