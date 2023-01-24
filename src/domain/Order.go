package domain

import "github.com/ZaphCode/clean-arch/src/utils"

//* Model

type Order struct {
	utils.DBModel
}

//* Service

type OrderService interface {
	utils.ServiceCrudOperations[Order]
}

//* Repository

type OrderRepository interface {
	utils.RepositoryCrudOperations[Order]
}
