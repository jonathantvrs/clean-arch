package handler

import (
	"context"
	"log"
	"net"
	"order-service/internal/service"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "order-service/proto"
)

type GRPCServer struct {
	Server *grpc.Server
}

func NewGRPCServer(orderService *service.OrderService) *GRPCServer {
	grpcServer := grpc.NewServer()

	pb.RegisterOrderServiceServer(grpcServer, &orderServiceServer{
		orderService: orderService,
	})

	reflection.Register(grpcServer)

	return &GRPCServer{
		Server: grpcServer,
	}
}

func (s *GRPCServer) Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	log.Printf("Starting gRPC server on port %s", port)
	return s.Server.Serve(lis)
}

func (s *GRPCServer) Shutdown() {
	s.Server.GracefulStop()
}

type orderServiceServer struct {
	pb.UnimplementedOrderServiceServer
	orderService *service.OrderService
}

func (s *orderServiceServer) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := s.orderService.ListOrders(ctx)
	if err != nil {
		return nil, err
	}

	var pbOrders []*pb.Order
	for _, o := range orders {
		pbOrders = append(pbOrders, &pb.Order{
			Id:          int32(o.ID),
			ProductName: o.ProductName,
			Quantity:    int32(o.Quantity),
			CreatedAt:   o.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.ListOrdersResponse{Orders: pbOrders}, nil
}
