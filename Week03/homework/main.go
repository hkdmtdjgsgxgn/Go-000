// 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
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
		// Hello world, the web server
		helloHandler := func(w http.ResponseWriter, req *http.Request) {
			io.WriteString(w, "Hello, world!\n")
		}

		http.HandleFunc("/hello", helloHandler)
		go func() {
			fmt.Println("Http Server start on 8081")
			serr <- s.ListenAndServe()
		}()
		select {
		case err := <-serr:
			err = s.Shutdown(context.Background())
			return err
		}
	})

	g.Go(func() error {
		signal.Notify(stop, syscall.SIGINT|syscall.SIGTERM|syscall.SIGKILL)
		<-stop
		select {
		case <-stop:
			return s.Shutdown(context.Background())
		}
	})
	if err := g.Wait(); err == nil {
		fmt.Println("http serve start.")
	} else {
		fmt.Printf("errgroup: %v", err)
	}
	close(stop)
	close(serr)
	fmt.Println(s.Close().Error())
}
