package domain

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID
	Username  string
	Email     string
	IsAdmin   bool
	Password  string
	ImageUrl  string
	Age       uint16
	CreatedAt int64
	UpdatedAt int64
}

type UserService interface {
	Create(*User) error
	GetAll() ([]User, error)
	GetByID(uuid.UUID) (*User, error)
	Delete(uuid.UUID) error
}

type UserRepository interface {
	Save(*User) error
	Find() ([]User, error)
	FindByID(uuid.UUID) (*User, error)
	Remove(uuid.UUID) error
	Update(uuid.UUID, *User) error
}
