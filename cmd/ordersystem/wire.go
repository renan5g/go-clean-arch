//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/renan5g/go-clean-arch/internal/application/repository"
	"github.com/renan5g/go-clean-arch/internal/application/usecase"
	"github.com/renan5g/go-clean-arch/internal/domain/event"
	"github.com/renan5g/go-clean-arch/internal/infra/database"
	"github.com/renan5g/go-clean-arch/pkg/events"
)

var provideOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(repository.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var provideOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

func MakeCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		provideOrderRepositoryDependency,
		provideOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

func MakeListOrdersUseCase(db *sql.DB) *usecase.ListOrdersUseCase {
	wire.Build(
		provideOrderRepositoryDependency,
		usecase.NewListOrdersUseCase,
	)
	return &usecase.ListOrdersUseCase{}
}
