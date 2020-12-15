package main

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		return errors.New("test")
	})

	err := g.Wait()
	fmt.Println(err)

	fmt.Println(ctx.Err())
}
