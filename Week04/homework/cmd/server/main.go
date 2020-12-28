package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/hi20160616/Go-000/Week04/homework/api/foobar/v1"
	"github.com/hi20160616/Go-000/Week04/homework/internal/pkg/fb_grpc"
	"github.com/hi20160616/Go-000/Week04/homework/internal/service"
	"golang.org/x/sync/errgroup"
)

const address = ":974"

func main() {
	// init service api
	fs := InitFoobarCase()
	service := service.NewFoobarService(fs)

	// register grpc service
	grpcServer := fb_grpc.NewServer(address)
	pb.RegisterFoobarServer(grpcServer, service)

	// context
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	// start grpc server
	g.Go(func() error {
		return grpcServer.Start(ctx)
	})

	// catch signals
	g.Go(func() error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-sigs:
			fmt.Println("")
			log.Printf("signal caught: %s, ready to quit...", sig.String())
			cancel()
		case <-ctx.Done():
			return ctx.Err()
		}
		return nil
	})

	// wait stop
	if err := g.Wait(); err != nil {
		log.Printf("error: %w", err)
	}
}
