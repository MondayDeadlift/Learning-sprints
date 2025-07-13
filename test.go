package main

import (
	"fmt"
	"sync"
)

// версия go 1.20
// что выведет код? как можно изменить вывод?
// что будет если убрать time.Sleep()? что вместо него?
func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			i := i
			fmt.Println(i)
		}()
	}
	wg.Wait()
}
