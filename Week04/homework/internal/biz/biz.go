package biz

import "math/rand"
import (
	"github.com/hi20160616/Go-000/Week04/homework/internal/data"
)

type Biz struct {
	id  int64
	key uint32
}

// MakeKey generate a random uint32 number with seed of `b.id`.
func (b *Biz) MakeKey(d data.Data) {
	r := rand.New(rand.NewSource(b.id))
	b.key = r.Uint32()
}
