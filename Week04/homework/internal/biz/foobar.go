package biz

type Foobar struct {
	ID  int32
	Foo string
	Bar int32
}

type FoobarRepo interface {
	Save(*Foobar) int32
}

type FoobarCase struct {
	repo FoobarRepo
}

func NewFoobarCase(repo FoobarRepo) *FoobarCase {
	return &FoobarCase{repo: repo}
}

func (f *FoobarCase) SaveFoobar(fb *Foobar) {
	id := f.repo.Save(fb)
	fb.ID = id
}
