package handler

import (
	"context"
	"log"
	"net/http"
	"order-service/internal/service"

	"order-service/internal/handler/generated"
	"order-service/internal/handler/model"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

type GraphQLServer struct {
	server *handler.Server
}

func NewGraphQLServer(orderService *service.OrderService) *GraphQLServer {
	resolver := &Resolver{orderService: orderService}
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: resolver,
			},
		),
	)

	return &GraphQLServer{
		server: srv,
	}
}

func (s *GraphQLServer) Start(port string) error {
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", s.server)

	log.Printf("Starting GraphQL server on port %s", port)
	return http.ListenAndServe(":"+port, nil)
}

type Resolver struct {
	orderService *service.OrderService
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r.orderService}
}

type queryResolver struct {
	orderService *service.OrderService
}

func (r *queryResolver) ListOrders(ctx context.Context) ([]*model.Order, error) {
	orders, err := r.orderService.ListOrders(ctx)
	if err != nil {
		return nil, err
	}

	var result []*model.Order
	for _, o := range orders {
		result = append(result, &model.Order{
			ID:          o.ID,
			ProductName: o.ProductName,
			Quantity:    o.Quantity,
			CreatedAt:   o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return result, nil
}
