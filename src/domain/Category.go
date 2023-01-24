package domain

import "github.com/ZaphCode/clean-arch/src/utils"

//* Model

type Category struct {
	utils.DBModel
}

//* Service

type CategoryService interface {
	utils.ServiceCrudOperations[Category]
}

//* Repository

type CategoryRepository interface {
	utils.RepositoryCrudOperations[Category]
}
