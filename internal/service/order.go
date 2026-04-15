package service

import (
	"context"
	"order-service/internal/model"
	"order-service/internal/repository"
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) ListOrders(ctx context.Context) ([]model.Order, error) {
	return s.repo.ListOrders(ctx)
}
