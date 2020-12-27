package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/hi20160616/Go-000/Week04/homework/api/hello/v1"
	"google.golang.org/grpc"
)

const (
	address     string = "localhost:8888"
	defaultName string = "world"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewFoobarClient(conn)

	// Contact server and print out response
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Say: %s", r.GetMessage())
}
