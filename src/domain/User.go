package domain

import (
	"github.com/google/uuid"
)

// * Model

type User struct {
	Model
	CustomerID    string `json:"customer_id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Role          string `json:"role"`
	Password      string `json:"password,omitempty"`
	VerifiedEmail bool   `json:"verified_email"`
	ImageUrl      string `json:"image_url"`
	Age           uint16 `json:"age"`
}

//* Service

type UserService interface {
	ServiceCrudOperations[User]
	GetByEmail(email string) (*User, error)
	GetByCredentials(email, pass string) (*User, error)
	VerifyEmail(ID uuid.UUID) error
	UpdatePassword(ID uuid.UUID, pass string) error
}

//* Repository

type UserRepository interface {
	RepositoryCrudOperations[User]
	FindByField(field string, val any) (*User, error)
	UpdateField(ID uuid.UUID, field string, val any) error
}
