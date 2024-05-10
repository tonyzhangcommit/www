package bussinesscode

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
	为sync.WaitGroup中Wait函数支持WaitTimeout功能( )
	要求sync.WaitGroup支持timeout功能
    如果timeout到了超时时间返回true
    如果WaitGroup自然结束返回false
*/

func WaitAddTOut() {
	wg := sync.WaitGroup{}
	ch := make(chan struct{}) // 这种定义方式，通常作为信号传输
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(num int, closech <-chan struct{}) {
			defer wg.Done()
			if _, ok := <-closech; ok {
				fmt.Println(num)
			}
		}(i, ch)
	}
	go func() {
		defer close(done)
		isclose := false
		go func() {
			for {
				select {
				case <-done:
					isclose = true
				}
			}
		}()
		for i := 0; i < 10; i++ {
			if !isclose {
				ch <- struct{}{}
				time.Sleep(time.Second * 3)
			}
		}
	}()
	// wg.Wait()
	if waitTimeOut(&wg, time.Second*10) {
		close(ch)
		fmt.Println("time out exit!")
		done <- true

	}
	time.Sleep(10 * time.Second)
}

func waitTimeOut(wg *sync.WaitGroup, timeout time.Duration) bool {
	timeOutCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	done := make(chan struct{})

	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
		return false
	case <-timeOutCtx.Done():
		return true
	}
}
