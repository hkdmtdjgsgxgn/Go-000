//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/hi20160616/Go-000/Week04/homework/internal/biz"
	"github.com/hi20160616/Go-000/Week04/homework/internal/data"
)

func InitFoobarCase() *biz.FoobarCase {
	wire.Build(biz.NewFoobarCase, data.NewFoobarRepo)
	return &biz.FoobarCase{}
}
