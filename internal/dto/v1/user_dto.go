package v1

import (
	"gin-user-management/internal/db/sqlc"
	"gin-user-management/internal/util"
)

type UserDTO struct {
	UUID      string `json:"uuid"`
	Name      string `json:"fullname"`
	Email     string `json:"email"`
	Age       *int   `json:"age"`
	Status    string `json:"status"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type CreateUserInput struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age" binding:"omitempty,gt=0"`
	Password string `json:"password" binding:"required,min=6"`
	Status   int    `json:"status" binding:"required,oneof=1 2 3"`
	Role     int    `json:"role" binding:"required,oneof=1 2 3"`
}

type UpdateUserInput struct {
	Name     string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age" binding:"omitempty,gt=0"`
	Password string `json:"password" binding:"omitempty,min=6"`
	Status   int    `json:"status" binding:"required,oneof=1 2 3"`
	Role     int    `json:"role" binding:"required,oneof=1 2 3"`
}

func (input *CreateUserInput) MapCreateInputToModel() sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		Email:    input.Email,
		Password: input.Password,
		Fullname: input.Fullname,
		Age:      util.ConvertToInt32Pointer(input.Age),
		Status:   int32(input.Status),
		Role:     int32(input.Role),
	}
}

func (input *UpdateUserInput) MapUpdateInputToModel() {
}

func MapUserToDTO(user sqlc.User) *UserDTO {
	dto := &UserDTO{
		UUID:      user.Uuid.String(),
		Name:      user.Fullname,
		Email:     user.Email,
		Role:      mapRoleText(int(user.Role)),
		Status:    mapStatusText(int(user.Status)),
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if user.Age != nil {
		age := int(*user.Age)
		dto.Age = &age
	}

	return dto
}

func mapStatusText(status int) string {
	switch status {
	case 1:
		return "Active"
	case 2:
		return "Inactive"
	case 3:
		return "Banned"
	default:
		return "None"
	}
}

func mapRoleText(status int) string {
	switch status {
	case 1:
		return "Admin"
	case 2:
		return "Moderator"
	case 3:
		return "Member"
	default:
		return "None"
	}
}
