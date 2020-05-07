package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/golang/protobuf/ptypes"
	v1 "github.com/julie-dev/go-grpc-http-rest-microservice-tutorial/pkg/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

const (
	apiVersion = "v1"
)

type GRPCClient struct {
	client v1.ToDoServiceClient
}

func getReminder() *timestamppb.Timestamp {
	t := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(t)

	return reminder
}

func (c *GRPCClient) Create (ctx context.Context) error {

	req := v1.CreateRequest{
		Api:  apiVersion,
		ToDo: &v1.ToDo{
			Title:       "title",
			Description: "description",
			Reminder:    getReminder(),
		},
	}
	res, err := c.client.Create(ctx, &req)
	if err != nil {
		return err
	}
	log.Printf("Create result: <%+v>\n\n", res)

	return nil
}

func (c *GRPCClient) ReadAll (ctx context.Context) error {

	req := v1.ReadAllRequest{
		Api:  apiVersion,
	}
	res, err := c.client.ReadAll(ctx, &req)
	if err != nil {
		return err
	}

	data, _ := json.MarshalIndent(res.ToDos, "", "  ")
	log.Printf("Create result: <%+v>\n\n", string(data))

	return nil
}

func main() {
	ctx := context.Background()

	address := flag.String("server", "", "gRPC Server in format host:port")
	flag.Parse()

	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	//
	var Client GRPCClient
	Client.client = v1.NewToDoServiceClient(conn)

	if err = Client.Create(ctx); err != nil {
		log.Fatalf("Create failed: %v", err)
	}

	if err = Client.ReadAll(ctx); err != nil {
		log.Fatalf("Read datas failed: %v", err)
	}

	return
}
