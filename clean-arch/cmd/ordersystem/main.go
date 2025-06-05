package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"time"

	"giovani-milanez/go-expert/clean-arch/configs"
	"giovani-milanez/go-expert/clean-arch/internal/event/handler"
	"giovani-milanez/go-expert/clean-arch/internal/infra/graph"
	"giovani-milanez/go-expert/clean-arch/internal/infra/grpc/pb"
	"giovani-milanez/go-expert/clean-arch/internal/infra/grpc/service"
	"giovani-milanez/go-expert/clean-arch/internal/infra/web/webserver"
	"giovani-milanez/go-expert/clean-arch/pkg/events"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"

	// migrations
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	fmt.Println("conectando ao banco...")
	// espera ate 30 segundos o banco subir
	db, err := getDbWithTimeout(30, configs)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	err = doMigration(db)
	if err != nil {
		panic(err)
	}

	rabbitMQChannel, err := getRabbitMQChannelWithTimeout(30, configs)
	if err != nil {
		panic(err)
	}

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUserCase  := NewListOrderUseCase(db)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("/order", webOrderHandler.Order)
	
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrderUserCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase: *listOrderUserCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", configs.GraphQLServerHost, configs.GraphQLServerPort), nil)
	if err != nil {
		panic(err)
	}
}

func getDbWithTimeout(tries int, configs *configs.Conf) (*sql.DB, error) {
	var db *sql.DB
	var err error
	count := 0
	for {
		if count >= tries {
			break
		}
		db, err = sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
		if err == nil {
			err = db.Ping()
		}
		if err == nil {
			break
		}
		fmt.Println(err.Error())
		time.Sleep(time.Second * 1)
		count++
	}
	return db, err
}

func doMigration(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
			"file://migrations",
			"mysql", 
			driver,
	)
	if err != nil {
		return err
	}
	
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func getRabbitMQChannelWithTimeout(tries int, cfg *configs.Conf) (*amqp.Channel, error) {
	var conn *amqp.Connection
	var err error
	count := 0
	for {
		if count >= tries {
			break
		}
		conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.RabbitMqUser, cfg.RabbitMqPassword, cfg.RabbitMqHost, cfg.RabbitMqPort))
		if err == nil {
			break
		}

		time.Sleep(time.Second * 1)
		count++
	}
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	_, err = ch.QueueDeclare("orders", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	err = ch.QueueBind("orders", "", "amq.direct", false, nil)
	if err != nil {
		return nil, err
	}

	return ch, nil
}
