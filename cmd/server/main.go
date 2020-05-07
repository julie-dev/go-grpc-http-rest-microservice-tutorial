package main

import (
	"fmt"
	"os"

	cmd "github.com/julie-dev/go-grpc-http-rest-microservice-tutorial/pkg/cmd/server"
)

func main() {
	if  err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
