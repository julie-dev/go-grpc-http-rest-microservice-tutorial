package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type HTTPClient struct {
	address *string
}

func (c *HTTPClient) Create(ctx context.Context) error {
	var body string

	reminder := getReminder()

	resp, err := http.Post(*c.address+"/v1/todo", "application/json",
		strings.NewReader(fmt.Sprintf(`
		{
			"api":"v1",
			"toDo": {
				"title":"title",
				"description":"description",
				"reminder":"%s"
			}
		}
	`, reminder)))
	if err != nil {
		return err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	body = string(bodyBytes)
	log.Printf("Create response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	return nil
}

func (c *HTTPClient) ReadAll(ctx context.Context) error {
	var body string

	resp, err := http.Get(*c.address + "/v1/todo/all")
	if err != nil {
		return err
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	body = string(bodyBytes)
	log.Printf("ReadAll response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	return nil
}

func (c *HTTPClient) Read(ctx context.Context, id string) error {
	var body string

	resp, err := http.Get(fmt.Sprintf("%s%s/%s", *c.address, "/v1/todo", id))
	if err != nil {
		log.Fatalf("failed to call Read method: %v", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	body = string(bodyBytes)
	log.Printf("Read response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	return nil
}

func getReminder() string {
	t := time.Now().In(time.UTC)
	return t.Format(time.RFC3339Nano)
}

func main() {
	ctx := context.Background()

	var Client HTTPClient
	Client.address = flag.String("server", "", "gRPC Server in format host:port")
	flag.Parse()

	var err error
	if err = Client.Create(ctx); err != nil {
		log.Fatalf("Create failed: %v", err)
	}

	if err = Client.ReadAll(ctx); err != nil {
		log.Fatalf("ReadAll failed: %v", err)
	}

	/*
		if err = Client.Read(ctx, "21"); err != nil {
			log.Fatalf("Read failed: %v", err)
		}
	*/
}
