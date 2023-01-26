package dtos

import "github.com/ZaphCode/clean-arch/src/domain"

type SigninDTO struct {
	Email    string `json:"email" validate:"required,email" example:"zaph@fapi.com"`
	Password string `json:"password" validate:"required" example:"menosfapi33"`
}

type SignupDTO struct {
	Username string `json:"username" validate:"required,min=4,max=15" example:"John doe"`
	Email    string `json:"email" validate:"required,email" example:"john@gmain.com"`
	Password string `json:"password" validate:"required,min=8" example:"password"`
	Age      uint16 `json:"age" validate:"required,number,gte=15" example:"18"`
}

func (dto SignupDTO) AdaptToUser() domain.User {
	return domain.User{
		Username: dto.Username,
		Email:    dto.Email,
		Password: dto.Password,
		Age:      dto.Age,
	}
}
