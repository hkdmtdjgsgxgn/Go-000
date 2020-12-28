package fb_grpc

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	address string
}

func NewServer(address string) *Server {
	srv := grpc.NewServer()
	return &Server{Server: srv, address: address}
}

func (s *Server) Start(ctx context.Context) error {
	l, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("grpc server start: %s", s.address)

	go func() {
		<-ctx.Done()
		s.GracefulStop()
		log.Println("grpc server graceful stopped.")
	}()

	return s.Serve(l)
}
