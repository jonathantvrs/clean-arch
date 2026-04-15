package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"order-service/internal/handler"
	"order-service/internal/repository"
	"order-service/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection
	db := connectDB()
	defer db.Close()

	// Setup dependencies
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)

	// Start servers
	restServer := startRESTServer(orderService)
	grpcServer := startGRPCServer(orderService)
	graphqlServer := startGraphQLServer(orderService)

	// Graceful shutdown
	waitForShutdown(restServer, grpcServer, graphqlServer)
}

func connectDB() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}

func startRESTServer(s *service.OrderService) *http.Server {
	router := mux.NewRouter()
	handler.SetupRoutes(router, s)

	srv := &http.Server{
		Addr:    ":" + os.Getenv("REST_PORT"),
		Handler: router,
	}

	go func() {
		log.Printf("Starting REST server on port %s", os.Getenv("REST_PORT"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("REST server error: %v", err)
		}
	}()

	return srv
}

func startGRPCServer(s *service.OrderService) *handler.GRPCServer {
	grpcServer := handler.NewGRPCServer(s)
	go func() {
		if err := grpcServer.Start(os.Getenv("GRPC_PORT")); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()
	return grpcServer
}

func startGraphQLServer(s *service.OrderService) *handler.GraphQLServer {
	graphqlServer := handler.NewGraphQLServer(s)
	go func() {
		if err := graphqlServer.Start(os.Getenv("GRAPHQL_PORT")); err != nil {
			log.Fatalf("GraphQL server error: %v", err)
		}
	}()
	return graphqlServer
}

func waitForShutdown(servers ...interface{}) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, s := range servers {
		switch server := s.(type) {
		case *http.Server:
			log.Println("Shutting down REST server...")
			server.Shutdown(ctx)
		case *handler.GRPCServer:
			log.Println("Shutting down gRPC server...")
			server.Shutdown()
		case *handler.GraphQLServer:
			log.Println("GraphQL server shutdown...")
		}
	}
}
