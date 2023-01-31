package dtos

import (
	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
)

type NewUserDTO struct {
	Username      string `json:"username" validate:"required,min=4,max=15" example:"John Doe"`
	Email         string `json:"email" validate:"required,email" example:"john@gmail.com"`
	VerifiedEmail bool   `json:"verified_email" example:"false"`
	Password      string `json:"password" validate:"required,min=8" example:"password123"`
	Role          string `json:"role" validate:"omitempty,oneof=user moderator" example:"user"`
	Age           uint16 `json:"age" validate:"required,number,gte=15" example:"20"`
	ImageUrl      string `json:"image_url" validate:"omitempty,url" example:"https://nwdistrict.ifas.ufl.edu/nat/files/2021/01/Groundhog.jpg"`
}

func (dto NewUserDTO) AdaptToUser() (user domain.User) {
	user.Username = dto.Username
	user.Email = dto.Email
	user.Role = dto.Role
	user.ImageUrl = dto.ImageUrl
	user.VerifiedEmail = dto.VerifiedEmail
	user.Password = dto.Password
	user.Age = dto.Age
	return
}

// ----------------------------------------------------------

type UserDTO struct { //? For documentation
	NewUserDTO
	CustomerID string    `json:"customer_id"`
	ID         uuid.UUID `json:"id" example:"8ded83fe-93c8-11ed-ab0f-d8bbc1a27048"`
	CreatedAt  int64     `json:"created_at" example:"1674405183"`
	UpdatedAt  int64     `json:"updated_at" example:"1674405181"`
}

// ----------------------------------------------------------

type UpdateUserDTO struct {
	Username      string `json:"username" validate:"omitempty,min=4,max=15" example:"John Doe"`
	ImageUrl      string `json:"image_url" validate:"omitempty,url" example:"https://nwdistrict.ifas.ufl.edu/nat/files/2021/01/Groundhog.jpg"`
	Age           uint16 `json:"age" validate:"omitempty,number,gte=15" example:"20"`
	VerifiedEmail *bool  `json:"verified_email" validate:"omitempty"`
	Role          string `json:"role" validate:"omitempty,oneof=user moderator" example:"user"`
}

func (dto UpdateUserDTO) AdaptToUser() (user domain.User) {
	user.Username = dto.Username
	user.ImageUrl = dto.ImageUrl
	if dto.VerifiedEmail != nil {
		user.VerifiedEmail = *dto.VerifiedEmail
	}
	user.Role = dto.Role
	user.Age = dto.Age
	return
}

// ----------------------------------------------------------
