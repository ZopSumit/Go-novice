package store

import (
	"layered/architecture/entities"
)

type Customer interface {
	GetByID(id int) (entities.Customer, error)
	GetByName(name string) ([]entities.Customer, error)
	Create(c entities.Customer) (entities.Customer, error)
}
