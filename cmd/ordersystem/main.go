package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/renan5g/go-clean-arch/config"
	"github.com/renan5g/go-clean-arch/internal/domain/event/handler"
	"github.com/renan5g/go-clean-arch/internal/infra/factory"
	"github.com/renan5g/go-clean-arch/internal/infra/graph"
	"github.com/renan5g/go-clean-arch/internal/infra/grpc/pb"
	"github.com/renan5g/go-clean-arch/internal/infra/grpc/service"
	"github.com/renan5g/go-clean-arch/pkg/events"
	"github.com/renan5g/go-clean-arch/pkg/rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	factory.MakeWebOrderHandle(db, r, eventDispatcher)
	fmt.Println("Starting Web server on port", configs.WebServerPort)
	go http.ListenAndServe(fmt.Sprintf(":%s", configs.WebServerPort), r)

	createOrderUseCase := factory.MakeCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := factory.MakeListOrdersUseCase(db)

	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase:   *listOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(fmt.Sprintf(":%s", configs.GraphQLServerPort), nil)
}
