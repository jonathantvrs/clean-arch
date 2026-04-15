package handler

import (
	"context"
	"order-service/internal/model"
)

type QueryResolver interface {
	ListOrders(ctx context.Context) ([]*model.Order, error)
}
