package main

import (
	"fmt"
	"strconv"
	"sync"
)

// что выведет код? если есть проблемы их нужно исправить
func main() {
	var wc sync.WaitGroup
	var wg sync.WaitGroup
	m := make(chan string, 3)
	for i := 0; i < 5; i++ {
		wc.Add(1)
		go func(mm chan<- string, i int) {
			defer wc.Done()
			mm <- fmt.Sprintf("Goroutine %s", strconv.Itoa(i))
		}(m, i)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			n, ok := <-m
			if !ok {
				return
			}
			fmt.Println(n)
		}
	}()

	wc.Wait()
	close(m)
	wg.Wait()
}
