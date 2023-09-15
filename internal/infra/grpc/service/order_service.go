package service

import (
	"context"

	"github.com/renan5g/go-clean-arch/internal/application/usecase"
	"github.com/renan5g/go-clean-arch/internal/infra/grpc/pb"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrdersUseCase  usecase.ListOrdersUseCase
}

func NewOrderService(
	createOrderUseCase usecase.CreateOrderUseCase,
	ListOrdersUseCase usecase.ListOrdersUseCase,
) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  ListOrdersUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.Order, error) {
	dto := usecase.OrderInputDTO{
		Price: in.Price,
		Tax:   in.Tax,
	}
	output, err := s.CreateOrderUseCase.Execute(&dto)
	if err != nil {
		return nil, err
	}
	return &pb.Order{
		Id:         output.ID,
		Price:      output.Price,
		Tax:        output.Tax,
		FinalPrice: output.FinalPrice,
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.Black) (*pb.OrderList, error) {
	orders, err := s.ListOrdersUseCase.Execute()
	if err != nil {
		return nil, err
	}

	var ordersOut []*pb.Order
	for _, o := range orders {
		order := &pb.Order{
			Id:         o.ID,
			Price:      o.Price,
			Tax:        o.Tax,
			FinalPrice: o.FinalPrice,
		}

		ordersOut = append(ordersOut, order)
	}

	return &pb.OrderList{Orders: ordersOut}, nil
}
