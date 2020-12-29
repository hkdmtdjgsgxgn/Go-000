package biz

type Greeter struct {
	ID   int32
	Name string
	Msg  string
}

type GreeterCase struct {
	repo GreeterRepo
}

type GreeterRepo interface {
	SayHi(*Greeter)
}

func (g *GreeterCase) SayHiThere(gg *Greeter) {
	g.repo.SayHi(gg)
}

func NewGreeterCase(repo GreeterRepo) *GreeterCase {
	return &GreeterCase{repo: repo}
}
