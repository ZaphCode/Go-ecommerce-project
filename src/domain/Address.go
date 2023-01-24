package domain

import "github.com/ZaphCode/clean-arch/src/utils"

//* Model

type Address struct {
	utils.DBModel
}

//* Service

type AddressService interface {
	utils.ServiceCrudOperations[Address]
}

//* Repository

type AddressRepository interface {
	utils.RepositoryCrudOperations[Address]
}
