package service

import (
	"context"

	v1 "github.com/hi20160616/Go-000/Week04/homework/api/foobar/v1"
	"github.com/hi20160616/Go-000/Week04/homework/internal/biz"
)

type FoobarService struct {
	fc *biz.FoobarCase
	v1.UnimplementedFoobarServer
}

func NewFoobarService(fc *biz.FoobarCase) v1.FoobarServer {
	return &FoobarService{fc: fc}
}

func (fs *FoobarService) RegisteFoobar(ctx context.Context, r *v1.FoobarRequest) (*v1.FoobarReply, error) {
	// dto -> do
	f := &biz.Foobar{Foo: r.Foo, Bar: r.Bar}

	// call biz
	fs.fc.SaveFoobar(f)

	//rt
	return &v1.FoobarReply{Id: f.ID}, nil
}
