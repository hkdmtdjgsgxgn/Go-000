package data

import (
	"log"

	"github.com/hi20160616/Go-000/Week04/helloworld/internal/biz"
)

var _ biz.GreeterRepo = new(greeterRepo)

type greeterRepo struct{}

func (g *greeterRepo) SayHi(biz.Greeter) {
	log.Printf("Hi there! data!")
}

func NewGreeterRepo() biz.GreeterRepo {
	return &greeterRepo{}
}
