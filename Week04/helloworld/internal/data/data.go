package data

import (
	"log"

	"github.com/hi20160616/Go-000/Week04/helloworld/internal/biz"
)

var _ biz.HelloRepo = new(helloRepo)

type helloRepo struct{}

const defaultID = 100

func NewHelloRepo() biz.HelloRepo {
	return &helloRepo{}
}

func (hr *helloRepo) GetID(h *biz.Hello) int32 {
	log.Printf("Got hello, name: %s, age: %d, id: %d", h.Name, h.Age, h.ID)
	return defaultID
}
