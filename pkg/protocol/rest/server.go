package rest

import (
	"context"
	v1 "github.com/julie-dev/go-grpc-http-rest-microservice-tutorial/pkg/api/v1"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

func RunServer(ctx context.Context, grpcPort, httpPort string) error {
	ctx, cancle := context.WithCancel(ctx)
	defer cancle()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := v1.RegisterToDoServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts); err != nil {
		log.Fatal("Failed to start HTTP gateway: %v", err)
	}

	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		log.Println("Shutting down gRPC gateway...")
		_ = srv.Shutdown(ctx)
		//<-ctx.Done()
	}()

	log.Println("Starting gRPC gateway...")
	return srv.ListenAndServe()
}
