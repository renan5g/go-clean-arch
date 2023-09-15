package main

import (
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/renan5g/go-clean-arch/config"
	"github.com/renan5g/go-clean-arch/internal/domain/event/handler"
	"github.com/renan5g/go-clean-arch/internal/infra/database"
	"github.com/renan5g/go-clean-arch/internal/infra/graph"
	"github.com/renan5g/go-clean-arch/internal/infra/grpc/pb"
	"github.com/renan5g/go-clean-arch/internal/infra/grpc/service"
	"github.com/renan5g/go-clean-arch/internal/infra/web"
	web_handler "github.com/renan5g/go-clean-arch/internal/infra/web/handler"

	"github.com/renan5g/go-clean-arch/pkg/events"
	"github.com/renan5g/go-clean-arch/pkg/rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs := config.LoadConfig(".")

	db := database.Open(configs)
	defer db.Close()

	rabbitMQChannel := rabbitmq.OpenChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{RabbitMQChannel: rabbitMQChannel})

	createOrderUseCase := MakeCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := MakeListOrdersUseCase(db)

	webserver := web.NewWebServer(configs.WebServerPort)
	webserver.SetupMiddleware()
	webOrderHandler := web_handler.NewWebOrderHandler(*createOrderUseCase, *listOrderUseCase)
	webserver.SetupRoutes(*webOrderHandler)
	fmt.Println("Starting Web server on port", configs.WebServerPort)
	go webserver.Start()

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
