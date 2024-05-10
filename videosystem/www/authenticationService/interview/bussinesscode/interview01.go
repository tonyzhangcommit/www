package bussinesscode

import (
	"context"
	"log"
	"time"
)

/*
写出以下逻辑，要求每秒钟调用一次proc并保证程序不退出(什么手写代码？？？ )
*/

func proc() {
	panic("ok")
}

func Callproc(ctx context.Context) {
	go func() {
		for {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Println(r)
					}
				}()
				proc()
			}()
			time.Sleep(time.Second)
		}
	}()

	select {
	case <-ctx.Done():
		// 外部终止
		return
	}
}
