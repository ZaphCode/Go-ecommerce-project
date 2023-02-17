package shared

import (
	"time"

	"github.com/ZaphCode/clean-arch/src/domain"
	"github.com/google/uuid"
)

// TODO: move this models to utils package

var m1 = &domain.ExampleModel{
	Model: domain.Model{
		ID:        uuid.MustParse("1551f9f0-825a-438c-9307-90cbc0bd5d63"),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	},
	Name:  "model 1",
	Tags:  []string{"A", "B", "C"},
	Check: true,
	Num:   143,
	Float: 42.5,
}

var m2 = &domain.ExampleModel{
	Model: domain.Model{
		ID:        uuid.MustParse("9f44a912-40f6-4ca6-b672-4911e3453443"),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	},
	Name:  "model 2",
	Check: true,
	Num:   143,
	Float: 42.5,
}

var m3 = &domain.ExampleModel{
	Model: domain.Model{
		ID:        uuid.MustParse("aa1a624e-555a-4b08-8bb4-3ed5aca074d7"),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	},
	Name:  "model 3",
	Num:   69,
	Float: 33.33,
}
