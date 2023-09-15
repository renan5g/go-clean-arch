package usecase

import (
	"github.com/renan5g/go-clean-arch/internal/application/repository"
	"github.com/renan5g/go-clean-arch/internal/domain/entity"
	"github.com/renan5g/go-clean-arch/pkg/events"
)

type OrderInputDTO struct {
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type CreateOrderUseCase struct {
	OrderRepo       repository.OrderRepositoryInterface
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewCreateOrderUseCase(
	repo repository.OrderRepositoryInterface,
	orderCreated events.EventInterface,
	eventDispatcher events.EventDispatcherInterface,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepo:       repo,
		OrderCreated:    orderCreated,
		EventDispatcher: eventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input *OrderInputDTO) (*OrderOutputDTO, error) {
	order, err := entity.NewOrder(input.Price, input.Tax)
	if err != nil {
		return nil, err
	}

	if err := order.CalculateFinalPrice(); err != nil {
		return nil, err
	}

	if err := c.OrderRepo.Save(order); err != nil {
		return nil, err
	}

	output := &OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}

	c.OrderCreated.SetPayload(output)
	c.EventDispatcher.Dispatch(c.OrderCreated)

	return output, nil
}
