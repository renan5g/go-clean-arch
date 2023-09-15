package factory

import (
	"database/sql"

	"github.com/renan5g/go-clean-arch/internal/application/usecase"
	"github.com/renan5g/go-clean-arch/internal/domain/event"
	"github.com/renan5g/go-clean-arch/internal/infra/database"
	"github.com/renan5g/go-clean-arch/pkg/events"
)

func MakeCreateOrderUseCase(
	db *sql.DB,
	eventDispatcher events.EventDispatcherInterface,
) *usecase.CreateOrderUseCase {
	repo := database.NewOrderRepository(db)
	orderCreated := event.NewOrderCreated()
	createOrderUseCase := usecase.NewCreateOrderUseCase(repo, orderCreated, eventDispatcher)
	return createOrderUseCase
}

func MakeListOrdersUseCase(db *sql.DB) *usecase.ListOrderUseCase {
	repo := database.NewOrderRepository(db)
	listOrderUseCase := usecase.NewListOrderUseCase(repo)
	return listOrderUseCase
}
