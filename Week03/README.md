# # Homework

>#学号: G20200607010680
>#班级: 1
>#作业链接:[https://github.com/hi20160616/Go-000/tree/main/Week03/homework](https://github.com/hi20160616/Go-000/tree/main/Week03/homework)

基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

以上作业，要求提交到 Github 上面，Week03 作业提交地址：

[https://github.com/Go-000/Go-000/issues/69 ](https://github.com/Go-000/Go-000/issues/69)

作业提交截止时间 12 月 9 日（周三）23:59 前。


# # TODO

- [ ] study and learn sync.atomic atomic.value
- [ ] study and learn Redis COW BGSave
- [ ] 阅读`errgroup`源码
- [ ] 阅读`sync.Pool`源码
- [ ] 阅读`sync.atomic`源码
- [ ] 阅读`context`源码
# # Goroutine

go func 的包装：

```go
func Go(x func()){
    go func(){
        defer func(){
            if err := recover(); err != nil {
                fmt.Println(err)
            }
        }()
        x()
    }()
}
```
这样，`Go(x)`来跑的`goroutine`就不再是`野生 goroutine`，不会因为panic而终止整个应用了。
## Processes and Threads

* 启动App，操作系统会为App创建一个Process（进程）
    * App 里可以有很多个Thread（线程）
* 线程是操作系统调度的一种执行路径，用于让处理器执行我们在函数中编写的代码。
    * 一个进程从一个线程开始——主线程
    * 主线程终止，则进程终止
    * 主线程是应用程序的原点
    * 主线程可以启动更多线程，这些线程又可以启动更多线程
* 线程属于哪个进程，在哪个可用处理器上运行，每个OS自有安排
## Concurrency and Parallelism

[https://www.jianshu.com/p/b11e251d3dc7](https://www.jianshu.com/p/b11e251d3dc7)

并发(concurrency):一个处理器同时处理多个任务。| 逻辑上的同时发生

并行(Parallelism): 多个处理器或多核处理器同时处理多个不同的任务。| 物理上的同时发生

>并发和并行的区别就是一个人同时吃三个馒头和三个人同时吃三个馒头。😆
## ![图片](https://uploader.shimo.im/f/vfBczkbn3NnhGh0s.JPG!thumbnail?fileGuid=yHjxc3dxVRcdqPxH)

1. Keep yourself busy or do the work yourself
2. Leave concurrency to the caller
* biz 的工作应该自己去做而不是委派给野生goroutine
* 一定是调用者去决定是否go func，是后台执行还是前台执行
* goroutine生命周期应该有你自己来管理，也就是说，你一定有手段来控制goroutine什么时候停止
* 如果委派给了某个goroutine要做到：Never start a goroutine without knowing when it will stop
* Any time you start a goroutine you must ask yourself:
    * When will it terminate?
    * What could prevent it from terminating?
* **Only use**`log.Fatal()`**from main.main or init functions.**

Example:

```go
func serve(addr string, handler http.Handler, stop <-chan struct{}) error {
  s := http.Server{
    Addr: addr,
    Handler: handler,
  }
  
  go func(){
    <-stop // wait for stop signal
    s.Shutdown(context.Background())
  }()
  
  return s.ListenAndServe()
}
func main(){
  done := make(chan error, 2) // because there are 2 go func() below
  stop := make(chan struct{})
  go func(){
    done <- serveDebug(stop)
  }()
  go func(){
    done <- serveApp(stop)
  }()
  
  var stopped bool
  for i := 0; i < cap(done); i++{
    if err := <-done; err != nil {
      fmt.Println("error: %v", err)
    }
    if !stopped{
      stopped = true
      close(stop)
    }
  }
}
```
上面的代码太经典了，背下来！
上面的代码都发生了什么：

* `serve`里的`go func()`是在Listen`s.ListenAndServe()`,当它退出的时候，会返回error，这个error会被main里的循环感知到，打印错误，并`close(stop)`从而实现优雅退出。当`serveApp(done)`执行完毕会给`done`传递一个`error`，这个`error`会被`if err := <-done`接收并打印，同时堵塞，指导第二个`done`到来，并打印。`if stopped`在第一次错误打印后变为true从而阻止第二次调用`close(stop)`。我只有若干卧槽来表达此时的惊叹，太TM巧妙了，这是谁TM想出来的，是不是毛先生！
* `close(stop)`可以给stop一个`0`signal，从而让`serve`里的`stop`阻塞取消，进而运行`s.Shutdown`从而关闭`s.ListenAndServe()`elegantly！
* `done`是知道什么时候推出，`stop`用来通知两个goroutine优雅退出
* `main`里的`for`执行2次是等待两个goroutine都退出了再退出main。
## *** Concurrency 3 core tech points

1. 把并发交给调用者，由调用者来控制goroutine
2. 搞清楚goroutine何时退出
3. 能够控制goroutine何时退出，管控它的生命周期

Example(01:11:42)[https://github.com/hi20160616/Go-000/blob/main/Week03/5.go](https://github.com/hi20160616/Go-000/blob/main/Week03/5.go)

```go
package main
import (
	"context"
	"fmt"
	"time"
)
func main() {
	tr := NewTracker()
	go tr.Run()
	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "test")
	_ = tr.Event(context.Background(), "test")
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()
	tr.Shutdown(ctx)
}
func NewTracker() *Tracker {
	return &Tracker{
		ch: make(chan string, 10),
	}
}
type Tracker struct {
	ch   chan string
	stop chan struct{}
}
func (t *Tracker) Event(ctx context.Context, data string) error {
	select {
	case t.ch <- data:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (t *Tracker) Run() {
	for data := range t.ch {
		time.Sleep(1 * time.Second)
		fmt.Println(data)
	}
	t.stop <- struct{}{}
}
func (t *Tracker) Shutdown(ctx context.Context) {
	close(t.ch)
	select {
	case <-t.stop:
	case <-ctx.Done():
	}
}
```
# # Memory model

* [https://golang.org/ref/mem](https://golang.org/ref/mem)
* Happen-Before
* 编译器和处理器只有在不会改变这个goroutine的行为时才可能修改读和写的执行顺序(01:35:00)
## Memory Reordering

CPU设计者为了榨干CPU性能，各种手段都上，比如流水线，分支预测等。其中为了提高读写内存的效率，会对读写指令进行重新排列，这就是内存重排(Memory Reordering)。

类似手段还有CPU重排和编译器重排。

## Memory model

* 如果事件e1发生在事件e2前，我们可以说e2发生在e1后。如果e1不发生在e2前也不发生在e2后，我们就说e1和e2是并发的
* 对变量v的零值初始化在内存模型中表现的与写操作相同。
* 对大于 single machine word 的变量的读写操作，表现的像是以不确定的顺序对多个 single   machine word 的变量操作。

[https://jianshu.com/p/5e44168f47a3](https://jianshu.com/p/5e44168f47a3)

# # Package sync

## Share Memory By Communicating

**Do not communicate by sharing memory; instead, share memory by communicating.**

* **尽量用channel（Go的通信）来解决同步访问内存的问题，而不是用互斥锁**
## Detecting Race Conditions With Go

8.go:

```go
// go build -o ./8/ -race ./8/8.go
// ./8/8
// ==================
// WARNING: DATA RACE
// Read at 0x000001279d08 by goroutine 10:
//   main.Routine()
//       /Users/foobar/go/src/github.com/foobar/Go-000/Week03/8/8.go:63 +0x47
//
// Previous write at 0x000001279d08 by goroutine 7:
//   main.Routine()
//       /Users/foobar/go/src/github.com/foobar/Go-000/Week03/8/8.go:66 +0x74
//
// Goroutine 10 (running) created at:
//   main.main()
//       /Users/foobar/go/src/github.com/foobar/Go-000/Week03/8/8.go:55 +0x72
//
// Goroutine 7 (finished) created at:
//   main.main()
//       /Users/foobar/go/src/github.com/foobar/Go-000/Week03/8/8.go:55 +0x72
// ==================
// ==================
// WARNING: DATA RACE
// Read at 0x000001279d08 by goroutine 8:
//   main.Routine()
//       /Users/foobar/go/src/github.com/foobar/Go-000/Week03/8/8.go:63 +0x47
//
// Previous write at 0x000001279d08 by goroutine 7:
//   main.Routine()
//       /Users/foobar/go/src/github.com/foobar/Go-000/Week03/8/8.go:66 +0x74
//
// Goroutine 8 (running) created at:
//   main.main()
//       /Users/foobar/go/src/github.com/foobar/Go-000/Week03/8/8.go:55 +0x72
//
// Goroutine 7 (finished) created at:
//   main.main()
//       /Users/foobar/go/src/github.com/foobar/Go-000/Week03/8/8.go:55 +0x72
// ==================
// Final Counter: 4
// Found 2 data race(s)
package main
import (
"fmt"
"sync"
"time"
)
var Wait sync.WaitGroup
var Counter int = 0
func main() {
for routine := 1; routine <=2; routine++ {
Wait.Add(1)
go Routine(routine)
}
Wait.Wait()
fmt.Printf("Final Counter: %d\n", Counter)
}
func Routine(id int) {
for count := 0; count < 2; count++ {
value := Counter
time.Sleep(1 * time.Nanosecond)
value++
Counter = value
}
Wait.Done()
}
```
1. Data race1: 63:`value := Counter`读取的Counter很可能被另外的goroutine修改了
2. Data race2: 66:`Counter = value`要写入的Counter很可能被另外的goroutine修改了
3. `time.Sleep()`会触发上下文切换
* **没有安全的data race(safe data race)**，您的程序要么没有  data race 要么其操作未定义
* Data race 的两个点：
    * 原子性
    * 可见性
* 锁里面的代码越简单越短越轻量越简单越好
## sync.atomic

* CopyOnWrite(COW)
* 走进程缓存，效率非常高
* 针对写小读多的场景来实现无锁访问共享数据。
## Mutex

* 针对写多读少的场景
* 下面的例子是**🔐锁饥饿**

![图片](https://uploader.shimo.im/f/ApkAYXZEVjFrRxMo.png!thumbnail?fileGuid=yHjxc3dxVRcdqPxH)

![图片](https://uploader.shimo.im/f/dw3tEs0gIxHNWAP4.png!thumbnail?fileGuid=yHjxc3dxVRcdqPxH)

    * 如果在实际业务场景中 goroutine 2 是关键逻辑，很可能，在非常饥饿的情况下，goroutine 2 就约等于不执行。
    * goroutine 1 总是抢到 mu，所以 goroutine 2 就只好总是 park to the goroutine scheduler。
* 几种 Mutex 的实现：
    * Barging: 吞吐量优先，不公平
    * Handsoff: 公平性优先，但是吞吐量低。完美平衡两个 goroutine，但是会降低性能。不牵手绝不放手。
    * Spinning: 自旋
## errgroup

* 并行工作流
* 错误处理或优雅降级
* context 传播和级联、取消
* 举报变量+闭包
* 毛先生踩到的坑，他们在kratos里自己重新封装了 errgroup
    * `errgroup.Go()`它产生Panic导致main退出
    * context作用域问题，`WithContext`的context作用域只在errgroup中
    * 传错context的问题

Doc:[https://pkg.go.dev/golang.org/x/sync/errgroup#Group.Go](https://pkg.go.dev/golang.org/x/sync/errgroup#Group.Go)

Eg:[https://github.com/hi20160616/Go-000/tree/main/](https://github.com/hi20160616/Go-000/tree/main/)[Week03/10/10.go](https://github.com/hi20160616/Go-000/tree/main/Week03/10/10.go)

eg:

```go
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
```
## sync.Pool

* 就是用来高频的内存申请
* 用来保存和复用临时对象以减少内存分配，降低GC压力
* Request-Driven特别合适
# # Package context

* Incoming requests to a server should create a Context.
    * 当一个请求进来，第一步是创建一个context. withtimeout, withcancel.
* Outgoing calls to servers should accept a Context.
    * 你调别人的时候一定要显式的传一个Context
* Do not store Contexts inside a struct type; instead, pass a Context explicitly to each function that needs it.
    * 显式传递: func(ctx context.Context, args...)，而不是包到 struct 里去传递
* The chain of function calls between them must propagate the Context.
    * 函数链被调用的过程中必须传播context
* Replace a Context using WithCancel, WithDeadline, WithTimeout, or WithValue.
* When a Context is canceled, all Contexts derived from it are also canceled.
* The same Context may be passed to functions running in different goroutines; Contexts are safe for simultaneous use by multiple goroutines.
* Do not pass a nil Context, even if a function permits it. Pass a TODO context if you are unsure about which Context to use.
    * 不要传nil context，不知道传啥的时候传个TODO
* Use context values only for request-scoped data that transits processes and APIs, not for passing optional parameters to functions.
    * 不要在Context里传一些带业务逻辑的数据，如果一定要传，尽量显式传参
* All blocking/long operations should be cancelable.
    * 所有耗时或block的一定要传递context，让它可以被取消
* Context.Value obscures your program’s flow.
    * Context.Value 它不应该影响业务逻辑和功能
* Context.Value should inform, not control.
    * 它只是个信号，只是一个具体的值，它不是一个控制逻辑，just for debug, log, route...
    * 不要从context里拿出个用户id出来，很恶心
* Try not to use context.Value.

[https://talks.golang.org/2014/gotham-context.slide#1](https://talks.golang.org/2014/gotham-context.slide#1)

# # Channels

## Unbuffered Channels

* 无缓冲信道的本质是保证同步
* Receuve 先于Send发生
* pros and cons：100% 保证能收到，但延时未知
## Buffered Channels

* Send 先于 Receive 发生
* pros and cons: 延迟更小，但是不保证数据到达，越大的buffer，越小的保证到达，buffer==1时，给你保证延迟一个消息到达
* 吞吐是靠多个goroutine 消费，buffer 的size只能影响延迟
## Go Concurrency Patterns

* Channel 一定是交给发送者来 close channel.

# # References

https://www.ardanlabs.com/blog/2018/11/goroutine-leaks-the-forgotten-sender.html
https://www.ardanlabs.com/blog/2019/04/concurrency-trap-2-incomplete-work.html
https://www.ardanlabs.com/blog/2014/01/concurrency-goroutines-and-gomaxprocs.html
https://dave.cheney.net/practical-go/presentations/qcon-china.html#_concurrency
https://golang.org/ref/mem
https://blog.csdn.net/caoshangpa/article/details/78853919
https://blog.csdn.net/qcrao/article/details/92759907
https://cch123.github.io/ooo/
https://blog.golang.org/codelab-share
https://dave.cheney.net/2018/01/06/if-aligned-memory-writes-are-atomic-why-do-we-need-the-sync-atomic-package
http://blog.golang.org/race-detector
https://dave.cheney.net/2014/06/27/ice-cream-makers-and-data-races
https://www.ardanlabs.com/blog/2014/06/ice-cream-makers-and-data-races-part-ii.html
https://medium.com/a-journey-with-go/go-how-to-reduce-lock-contention-with-the-atomic-package-ba3b2664b549
https://medium.com/a-journey-with-go/go-discovery-of-the-trace-package-e5a821743c3c
https://medium.com/a-journey-with-go/go-mutex-and-starvation-3f4f4e75ad50
https://www.ardanlabs.com/blog/2017/10/the-behavior-of-channels.html
https://medium.com/a-journey-with-go/go-buffered-and-unbuffered-channels-29a107c00268
https://medium.com/a-journey-with-go/go-ordering-in-select-statements-fd0ff80fd8d6
https://www.ardanlabs.com/blog/2017/10/the-behavior-of-channels.html
https://www.ardanlabs.com/blog/2014/02/the-nature-of-channels-in-go.html
https://www.ardanlabs.com/blog/2013/10/my-channel-select-bug.html
https://blog.golang.org/io2013-talk-concurrency
https://blog.golang.org/waza-talk
https://blog.golang.org/io2012-videos
https://blog.golang.org/concurrency-timeouts
https://blog.golang.org/pipelines
https://www.ardanlabs.com/blog/2014/02/running-queries-concurrently-against.html
https://blogtitle.github.io/go-advanced-concurrency-patterns-part-3-channels/
https://www.ardanlabs.com/blog/2013/05/thread-pooling-in-go-programming.html
https://www.ardanlabs.com/blog/2013/09/pool-go-routines-to-process-task.html
https://blogtitle.github.io/categories/concurrency/
https://medium.com/a-journey-with-go/go-context-and-cancellation-by-propagation-7a808bbc889c
https://blog.golang.org/context
https://www.ardanlabs.com/blog/2019/09/context-package-semantics-in-go.html
https://golang.org/ref/spec#Channel_types
https://drive.google.com/file/d/1nPdvhB0PutEJzdCq5ms6UI58dp50fcAN/view
https://medium.com/a-journey-with-go/go-context-and-cancellation-by-propagation-7a808bbc889c
https://blog.golang.org/context
https://www.ardanlabs.com/blog/2019/09/context-package-semantics-in-go.html
https://golang.org/doc/effective_go.html#concurrency
https://zhuanlan.zhihu.com/p/34417106?hmsr=toutiao.io
https://talks.golang.org/2014/gotham-context.slide#1
https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39
