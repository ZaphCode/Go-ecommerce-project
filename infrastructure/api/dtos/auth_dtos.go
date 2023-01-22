package dtos

import "github.com/ZaphCode/clean-arch/domain"

type SigninDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignupDTO struct {
	Username string `json:"username" validate:"required,min=4,max=15"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Age      uint16 `json:"age" validate:"required,number,min=15"`
}

func (dto SignupDTO) AdaptToUser() domain.User {
	return domain.User{
		Username: dto.Username,
		Email:    dto.Email,
		Password: dto.Password,
		Age:      dto.Age,
	}
}
