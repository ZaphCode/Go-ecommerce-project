package domain

import "github.com/ZaphCode/clean-arch/src/domain/shared"

//* Model

type PaymentMethod struct {
	shared.DomainModel
}

//* Service

type PaymentMethodService interface {
	shared.ServiceCrudOperations[PaymentMethod]
}

//* Repository

type PaymentMethodRepository interface {
	shared.RepositoryCrudOperations[PaymentMethod]
}
