// go build -o ./8/ -race ./8/8.go
// ./8/8
// ==================
// WARNING: DATA RACE
// Read at 0x000001279d08 by goroutine 10:
//   main.Routine()
//       /Users/wedojava/go/src/github.com/wedojava/Go-000/Week03/8/8.go:63 +0x47
//
// Previous write at 0x000001279d08 by goroutine 7:
//   main.Routine()
//       /Users/wedojava/go/src/github.com/wedojava/Go-000/Week03/8/8.go:66 +0x74
//
// Goroutine 10 (running) created at:
//   main.main()
//       /Users/wedojava/go/src/github.com/wedojava/Go-000/Week03/8/8.go:55 +0x72
//
// Goroutine 7 (finished) created at:
//   main.main()
//       /Users/wedojava/go/src/github.com/wedojava/Go-000/Week03/8/8.go:55 +0x72
// ==================
// ==================
// WARNING: DATA RACE
// Read at 0x000001279d08 by goroutine 8:
//   main.Routine()
//       /Users/wedojava/go/src/github.com/wedojava/Go-000/Week03/8/8.go:63 +0x47
//
// Previous write at 0x000001279d08 by goroutine 7:
//   main.Routine()
//       /Users/wedojava/go/src/github.com/wedojava/Go-000/Week03/8/8.go:66 +0x74
//
// Goroutine 8 (running) created at:
//   main.main()
//       /Users/wedojava/go/src/github.com/wedojava/Go-000/Week03/8/8.go:55 +0x72
//
// Goroutine 7 (finished) created at:
//   main.main()
//       /Users/wedojava/go/src/github.com/wedojava/Go-000/Week03/8/8.go:55 +0x72
// ==================
// Final Counter: 4
// Found 2 data race(s)
package main

import (
	"fmt"
	"sync"
	"time"
)

var Wait sync.WaitGroup
var Counter int = 0

func main() {
	for routine := 1; routine < 3; routine++ {
		Wait.Add(1)
		go Routine(routine)
	}
	Wait.Wait()
	fmt.Printf("Final Counter: %d\n", Counter)
}

func Routine(id int) {
	for count := 0; count < 2; count++ {
		value := Counter
		time.Sleep(1 * time.Nanosecond)
		value++
		Counter = value
	}
	Wait.Done()
}
