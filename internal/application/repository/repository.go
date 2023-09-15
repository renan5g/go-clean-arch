package repository

import "github.com/renan5g/go-clean-arch/internal/domain/entity"

type OrderRepositoryInterface interface {
	Save(order *entity.Order) error
	List() ([]*entity.Order, error)
	GetTotal() (int, error)
}
