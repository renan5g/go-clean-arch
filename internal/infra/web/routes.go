package web

import (
	"github.com/renan5g/go-clean-arch/internal/infra/web/handler"
)

func (s *WebServer) SetupRoutes(orderHandler handler.WebOrderHandler) {
	s.Router.Get("/orders", orderHandler.List)
	s.Router.Post("/orders", orderHandler.Create)
}
