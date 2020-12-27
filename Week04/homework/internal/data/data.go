package data

import (
	"fmt"
)

type Data struct {
}

func (d *Data) SaveIt(id, key uint32) {
	fmt.Printf("Save %v: %v, done.\n", id, key)
}
