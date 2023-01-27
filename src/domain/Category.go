package domain

import "github.com/ZaphCode/clean-arch/src/domain/shared"

//* Model

type Category struct {
	shared.DomainModel
}

//* Service

type CategoryService interface {
	shared.ServiceCrudOperations[Category]
}

//* Repository

type CategoryRepository interface {
	shared.RepositoryCrudOperations[Category]
}
