package usecase

import (
	"github.com/renan5g/go-clean-arch/internal/application/repository"
)

type OrderListOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type ListOrdersUseCase struct {
	OrderRepo repository.OrderRepositoryInterface
}

func NewListOrdersUseCase(repo repository.OrderRepositoryInterface) *ListOrdersUseCase {
	return &ListOrdersUseCase{OrderRepo: repo}
}

func (uc *ListOrdersUseCase) Execute() ([]*OrderListOutputDTO, error) {
	orders, err := uc.OrderRepo.List()
	if err != nil {
		return nil, err
	}

	var ordersOut []*OrderListOutputDTO
	for _, o := range orders {
		order := &OrderListOutputDTO{
			ID:         o.ID,
			Price:      o.Price,
			Tax:        o.Tax,
			FinalPrice: o.FinalPrice,
		}

		ordersOut = append(ordersOut, order)
	}

	return ordersOut, nil
}
