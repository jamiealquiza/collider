package main

import (
	"fmt"
	"time"
	"os"
	"sync"

	"github.com/jamiealquiza/collider"
	"github.com/jamiealquiza/tachymeter"
)

type Thing struct {
	value int
}

const (
	items = 1024
	readers = 128
	reads = 50000
)

func main() {
	c, err := collider.New(items)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i := 0; i < items; i++ {
		c.Add(&Thing{value: 0})
	}

	t := tachymeter.New(&tachymeter.Config{Size: readers*reads})
	wg := sync.WaitGroup{}
	wg.Add(readers)

	start := time.Now()
	for i := 0; i < readers; i++ {
		go reader(c, t, &wg)
	}
	t.SetWallTime(time.Since(start))

	wg.Wait()
	t.Calc().Dump()
}

func reader(c *collider.Ring, t *tachymeter.Tachymeter, wg *sync.WaitGroup) {
	for i := 0; i < reads; i++ {
		start := time.Now()
		_ = c.Get().(*Thing).value
		t.AddTime(time.Since(start))
	}

	wg.Done()
}
