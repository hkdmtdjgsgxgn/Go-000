package main

import (
	"context"
	"log"
	"net"

	pb "github.com/hi20160616/Go-000/Week04/homework/api/hello/v1"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFoobarServer
}

// SayHello implements hello.HelloResponse
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterFoobarServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
