package domain

import "github.com/ZaphCode/clean-arch/src/domain/shared"

//* Model

type Order struct {
	shared.DomainModel
}

//* Service

type OrderService interface {
	shared.ServiceCrudOperations[Order]
}

//* Repository

type OrderRepository interface {
	shared.RepositoryCrudOperations[Order]
}
