package data

import (
	"fmt"
	"github.com/hi20160616/Go-000/Week04/homework/internal/biz"
)

// error occur on compile while foobarRepo does not implemented biz.FoobarRepo
var _ biz.FoobarRepo = new(foobarRepo)

const foobarId = 100

type foobarRepo struct{}

func NewFoobarRepo() biz.FoobarRepo {
	return &foobarRepo{}
}

func (fr *foobarRepo) Save(f *biz.Foobar) int32 {
	fmt.Printf("save foobar, foo: %s, bar: %d, id: %d", f.Foo, f.Bar, foobarId)
	return foobarId
}
