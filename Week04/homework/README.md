## Week04 作业题目

> 按照自己的构想，写一个项目，满足基本的目录结构和工程， 代码需要包含对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动， 信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架。

reference repo: https://github.com/AngelovLee/Go-001/tree/main/Week04

## Implements

### Prerequisites

#### gRPC

https://grpc.io/docs/languages/go/quickstart/#prerequisites

```
brew install protobuf
export GO111MODULE=on  # Enable module mode
go get google.golang.org/protobuf/cmd/protoc-gen-go \
         google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go
```

`export PATH="$PATH:$(go env GOPATH)/bin"` append to `.zshrc`

After prerequistes prepared, coding and run below:

```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/foobar/v1/foobar.proto
```

#### wire

```
go get github.com/google/wire/cmd/wire
```
