package factory

import (
	"database/sql"

	"github.com/go-chi/chi"
	"github.com/renan5g/go-clean-arch/internal/application/usecase"
	"github.com/renan5g/go-clean-arch/internal/domain/event"
	"github.com/renan5g/go-clean-arch/internal/infra/database"
	"github.com/renan5g/go-clean-arch/internal/infra/web"
	"github.com/renan5g/go-clean-arch/pkg/events"
)

func MakeWebOrderHandle(
	db *sql.DB,
	r chi.Router,
	eventDispatcher events.EventDispatcherInterface,
) *web.WebOrderHandler {
	repo := database.NewOrderRepository(db)
	orderCreated := event.NewOrderCreated()
	createOrderUseCase := usecase.NewCreateOrderUseCase(repo, orderCreated, eventDispatcher)
	listOrderUseCase := usecase.NewListOrderUseCase(repo)
	handle := web.NewWebOrderHandler(r, *createOrderUseCase, *listOrderUseCase)
	return handle
}
