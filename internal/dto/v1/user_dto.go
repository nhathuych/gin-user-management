package v1

import (
	"gin-user-management/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserDTO struct {
	UUID      string `json:"uuid"`
	Name      string `json:"fullname"`
	Email     string `json:"email"`
	Age       *int32 `json:"age"`
	Status    string `json:"status"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type CreateUserInput struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Age      *int32 `json:"age" binding:"omitempty,gt=0"`
	Password string `json:"password" binding:"required,min=6"`
	Status   int32  `json:"status" binding:"required,oneof=1 2 3"`
	Role     int32  `json:"role" binding:"required,oneof=1 2 3"`
}

type UpdateUserInput struct {
	Fullname *string `json:"fullname" binding:"omitempty"`
	Age      *int32  `json:"age" binding:"omitempty,gt=0"`
	Password *string `json:"password" binding:"omitempty,min=6"`
	Status   *int32  `json:"status" binding:"omitempty,oneof=1 2 3"`
	Role     *int32  `json:"role" binding:"omitempty,oneof=1 2 3"`
}

type GetUserByUuidParam struct {
	Uuid string `uri:"uuid" binding:"uuid"`
}

type GetUsersParams struct {
	Search string `form:"search" binding:"omitempty,min=3,max=50,search"`
	Page   int32  `form:"page" binding:"omitempty,gte=1"`
	Limit  int32  `form:"limit" binding:"omitempty,gte=1,lte=1000"`
	Order  string `form:"order_by" binding:"omitempty,oneof=id"`
	Sort   string `form:"sort" binding:"omitempty,oneof=asc desc"`
}

func (input *CreateUserInput) MapCreateInputToModel() sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		Email:    input.Email,
		Password: input.Password,
		Fullname: input.Fullname,
		Age:      input.Age,
		Status:   input.Status,
		Role:     input.Role,
	}
}

func (input *UpdateUserInput) MapUpdateInputToModel(uuid uuid.UUID) sqlc.UpdateUserParams {
	return sqlc.UpdateUserParams{
		Uuid:     uuid,
		Password: input.Password,
		Fullname: input.Fullname,
		Age:      input.Age,
		Status:   input.Status,
		Role:     input.Role,
	}
}

func MapUserToDTO(user sqlc.User) *UserDTO {
	dto := &UserDTO{
		UUID:      user.Uuid.String(),
		Name:      user.Fullname,
		Email:     user.Email,
		Role:      mapRoleText(user.Role),
		Status:    mapStatusText(user.Status),
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if user.Age != nil {
		age := *user.Age
		dto.Age = &age
	}

	return dto
}

func MapUsersToDTOs(users []sqlc.User) []UserDTO {
	dtos := make([]UserDTO, 0, len(users))
	for _, user := range users {
		dtos = append(dtos, *MapUserToDTO(user))
	}

	return dtos
}

func mapStatusText(status int32) string {
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

func mapRoleText(status int32) string {
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
