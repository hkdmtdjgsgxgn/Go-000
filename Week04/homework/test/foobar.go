package main

import (
	"context"
	"log"
	"time"

	v1 "github.com/hi20160616/Go-000/Week04/homework/api/foobar/v1"
	"google.golang.org/grpc"
)

const (
	address = "localhost:974"
	foo     = "broadcast"
	bar     = 974
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("cannot connect: %v", err)
	}
	defer conn.Close()

	// grpc client
	c := v1.NewFoobarClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// call rpc
	r, err := c.RegisteFoobar(ctx, &v1.FoobarRequest{Foo: foo, Bar: bar})
	if err != nil {
		log.Fatalf("register foobar failed: %v", err)
	}
	log.Printf("register foobar id: %d", r.GetId())
}
