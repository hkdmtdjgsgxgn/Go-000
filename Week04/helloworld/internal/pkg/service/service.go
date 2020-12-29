package service

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
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	log.Printf("\ngrpc server start at: %s", s.address)

	go func() {
		<-ctx.Done()
		s.GracefulStop()
		log.Printf("grpc server gracefully stopped.")
	}()
	return s.Serve(lis)
}
