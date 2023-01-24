package dtos

import "github.com/ZaphCode/clean-arch/src/domain"

type NewUserDTO struct {
	Username      string `json:"username" validate:"required,min=4,max=15"`
	Email         string `json:"email" validate:"required,email"`
	Role          string `json:"role" validate:"omitempty,oneof=user moderator"`
	ImageUrl      string `json:"image_url" validate:"omitempty,url"`
	VerifiedEmail bool   `json:"verified_email"`
	Password      string `json:"password" validate:"required,min=8"`
	Age           uint16 `json:"age" validate:"required,number,gte=15"`
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

type UpdateUserDTO struct {
	Username      string `json:"username" validate:"required,min=4,max=15"`
	ImageUrl      string `json:"image_url" validate:"omitempty,url"`
	Age           uint16 `json:"age" validate:"required,number,gte=15"`
	VerifiedEmail *bool  `json:"verified_email" validate:"omitempty"`
	Role          string `json:"role" validate:"omitempty,oneof=user moderator"`
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
