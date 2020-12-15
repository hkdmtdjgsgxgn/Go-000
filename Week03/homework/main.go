// 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
// Reference: The best practice in my opion : https://github.com/XYZ0901/Go-000/blob/main/Week03/work/main.go
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	stop := make(chan os.Signal, 1)
	serr := make(chan error, 1)
	s := &http.Server{Addr: ":8081"}
	g := new(errgroup.Group)
	g.Go(func() error {
		go func() {
			serr <- s.ListenAndServe()
		}()
		helloHandler := func(w http.ResponseWriter, req *http.Request) {
			io.WriteString(w, "Hello, world!\n")
		}
		http.HandleFunc("/", helloHandler)
		select {
		case err := <-serr:
			close(serr)
			close(stop) // without close, it'll block at `a := <-stop`
			return err
		}
	})

	// server closed here
	g.Go(func() error {
		signal.Notify(stop, syscall.SIGINT|syscall.SIGTERM|syscall.SIGKILL)
		if a := <-stop; a != nil {
			fmt.Println("Got signal:", a)
		}
		return s.Shutdown(context.Background())
	})
	if err := g.Wait(); err == nil {
		fmt.Println("working...")
	} else {
		fmt.Printf("errgroup: %v\n", err)
	}
}
