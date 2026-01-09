package v1

import (
	"database/sql"
	"errors"
	"gin-user-management/internal/db/sqlc"
	"gin-user-management/internal/repository"
	"gin-user-management/internal/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetAll() {
	us.repo.GetAll()
}

func (us *userService) Create(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return sqlc.User{}, util.WrapError(err, "Failed to hash password.", util.ErrCodeInternal)
	}

	input.Password = string(hashedPassword)
	input.Email = util.NormalizeString(input.Email)

	user, err := us.repo.Create(context, input)
	if err != nil {
		return sqlc.User{}, util.WrapError(err, "Failed to create a new user.", util.ErrCodeInternal)
	}

	return user, nil
}

func (us *userService) GetByUUID() {
	us.repo.GetByUUID()
}

func (us *userService) Update(ctx *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	if input.Password != nil && *input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return sqlc.User{}, util.WrapError(err, "Failed to hash password.", util.ErrCodeInternal)
		}

		strPassword := string(hashedPassword)
		input.Password = &strPassword
	}

	user, err := us.repo.Update(context, input)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, util.NewError("User not found.", util.ErrCodeNotFound)
		}
		return sqlc.User{}, util.WrapError(err, "Failed to update user.", util.ErrCodeInternal)
	}

	return user, nil
}

func (us *userService) Delete() {
	us.repo.Delete()
}
