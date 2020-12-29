package biz

type Hello struct {
	ID   int32
	Name string
	Age  int32
}

type HelloRepo interface {
	GetID(*Hello) int32
}

type HelloCase struct {
	repo HelloRepo
}

func NewHelloCase(repo HelloRepo) *HelloCase {
	return &HelloCase{repo: repo}
}

func (hc *HelloCase) GetID(h *Hello) {
	id := hc.repo.GetID(h)
	h.ID = id
}
