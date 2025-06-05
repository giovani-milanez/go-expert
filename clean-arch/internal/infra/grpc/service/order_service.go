package service

import (
	"context"

	"giovani-milanez/go-expert/clean-arch/internal/infra/grpc/pb"
	"giovani-milanez/go-expert/clean-arch/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase usecase.CreateOrderUseCase
	ListOrderUseCase usecase.ListOrderUseCase
}

func NewOrderService(createOrderUseCase usecase.CreateOrderUseCase, listOrderUseCase usecase.ListOrderUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrderUseCase: listOrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.CreateOrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrder(context.Context, *pb.Blank) (*pb.ListOrderResponse, error) {
	output, err := s.ListOrderUseCase.Execute()
	if err != nil {
		return nil, err
	}

	resp := &pb.ListOrderResponse{}
	resp.Orders = []*pb.CreateOrderResponse{}

	for _, dto := range output {
		order := &pb.CreateOrderResponse{
			Id:         dto.ID,
			Price:      float32(dto.Price),
			Tax:        float32(dto.Tax),
			FinalPrice: float32(dto.FinalPrice),
		}
		resp.Orders = append(resp.Orders, order)
	}
	return resp, nil
}