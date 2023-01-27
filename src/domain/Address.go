package domain

import (
	"github.com/ZaphCode/clean-arch/src/domain/shared"
)

//* Model

type Address struct {
	shared.DomainModel
}

//* Service

type AddressService interface {
	shared.ServiceCrudOperations[Address]
}

//* Repository

type AddressRepository interface {
	shared.RepositoryCrudOperations[Address]
}
