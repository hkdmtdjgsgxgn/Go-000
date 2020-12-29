package service

import (
	"context"
	"log"

	pb "github.com/hi20160616/Go-000/Week04/helloworld/api/helloworld/v1"
)

const port = ":50051"

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "hello" + in.GetName()}, nil
}
