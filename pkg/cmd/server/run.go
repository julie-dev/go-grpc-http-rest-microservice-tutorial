package server

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/julie-dev/go-grpc-http-rest-microservice-tutorial/pkg/protocol/rest"
	v1 "github.com/julie-dev/go-grpc-http-rest-microservice-tutorial/pkg/service/v1"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julie-dev/go-grpc-http-rest-microservice-tutorial/pkg/protocol/grpc"
)

type Config struct {
	GRPCPort            string
	HTTPPort            string
	DatastoreDBHost     string
	DatastoreDBUser     string
	DatastoreDBPassword string
	DatastoreDBSchema   string
}

func check(port string) bool {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		return false
	}

	if portNum < 1 || portNum > 65535 {
		return false
	}

	return true
}

func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.HTTPPort, "http-port", "", "HTTP port to bind")
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "", "Database schema")
	flag.Parse()

	if check(cfg.GRPCPort) == false {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	if check(cfg.HTTPPort) == false {
		return fmt.Errorf("invalid TCP port for HTTP server: '%s'", cfg.HTTPPort)
	}

	param := "parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
		cfg.DatastoreDBHost,
		cfg.DatastoreDBSchema,
		param)

	// db is mysql handler
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to opendatabase: %v", err)
	}
	defer db.Close()

	// v1API is implement of ToDoServiceServer interface
	v1API := v1.NewTodoServiceServer(db)

	// run gRPC gateway server
	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	// Send gRPC configuration for running gRPC Server
	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
