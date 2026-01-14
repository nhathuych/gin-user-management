package v1

import (
	"database/sql"
	"errors"
	"gin-user-management/internal/db/sqlc"
	"gin-user-management/internal/repository"
	"gin-user-management/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (us *userService) GetAll(ctx *gin.Context, search, orderBy, sort string, page, limit int32, deleted bool) ([]sqlc.User, int32, error) {
	context := ctx.Request.Context()

	if sort == "" {
		sort = "desc"
	}
	if orderBy == "" {
		orderBy = "id"
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, err := us.repo.GetAllV2(context, search, orderBy, sort, limit, offset, deleted)
	if err != nil {
		return []sqlc.User{}, 0, util.WrapError(err, "Failed to retrieve users.", util.ErrCodeInternal)
	}

	total, err := us.repo.CountUsers(context, search, deleted)
	if err != nil {
		return []sqlc.User{}, 0, util.WrapError(err, "Failed to count users.", util.ErrCodeInternal)
	}

	return users, int32(total), nil
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

func (us *userService) GetByUUID(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()

	user, err := us.repo.GetByUUID(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, util.NewError("User not found.", util.ErrCodeNotFound)
		}
		return sqlc.User{}, util.WrapError(err, "Failed to retrieve user.", util.ErrCodeInternal)
	}

	return user, nil
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

func (us *userService) SoftDeleteUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()

	user, err := us.repo.SoftDeleteUser(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, util.NewError("User not found.", util.ErrCodeNotFound)
		}
		return sqlc.User{}, util.WrapError(err, "Failed to delete user.", util.ErrCodeInternal)
	}

	return user, nil
}

func (us *userService) RestoreUser(ctx *gin.Context, uuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()

	user, err := us.repo.RestoreUser(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, util.NewError("User not found or not deleted yet.", util.ErrCodeNotFound)
		}
		return sqlc.User{}, util.WrapError(err, "Failed to restore user.", util.ErrCodeInternal)
	}

	return user, nil
}

func (us *userService) HardDeleteUser(ctx *gin.Context, uuid uuid.UUID) error {
	context := ctx.Request.Context()

	_, err := us.repo.HardDeleteUser(context, uuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return util.NewError("User not found or not eligible for hard delete.", util.ErrCodeNotFound)
		}
		return util.WrapError(err, "Failed to delete user.", util.ErrCodeInternal)
	}

	return nil
}
