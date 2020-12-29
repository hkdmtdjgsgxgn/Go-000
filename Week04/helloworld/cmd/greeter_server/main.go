package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/hi20160616/Go-000/Week04/helloworld/api/helloworld/v1"
	"github.com/hi20160616/Go-000/Week04/helloworld/internal/pkg/service"
	"golang.org/x/sync/errgroup"
)

const (
	address = ":50051"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello again " + in.GetName()}, nil
}

func main() {
	// lis, err := net.Listen("tcp", port)
	// if err != nil {
	//         log.Fatalf("failed to listen: %v", err)
	// }
	// s := grpc.NewServer()
	// pb.RegisterGreeterServer(s, &server{})
	// if err := s.Serve(lis); err != nil {
	//         log.Fatalf("failed to serve: %v", err)
	// }

	s := service.NewServer(address)
	pb.RegisterGreeterServer(s, &server{})

	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error { return s.Start(ctx) })

	g.Go(func() error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-sigs:
			fmt.Println()
			log.Printf("signal caught: %s, reday to quit...", sig.String())
			cancel()
		case <-ctx.Done():
			return ctx.Err()
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		log.Printf("greeter_server main error: %v", err)
	}
}
