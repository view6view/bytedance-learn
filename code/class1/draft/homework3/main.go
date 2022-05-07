package main

import (
	"context"
	"fmt"
	"homework3/baidu"
	"homework3/caiyun"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, `usage: simpleDict WORD
	example: simpleDict hello
			`)
		os.Exit(1)
	}
	word := os.Args[1]

	// 创建上下文对象
	rep := make(chan string)
	// 关闭chan
	defer close(rep)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	go baidu.Query(word, rep)
	go caiyun.Query(word, rep)

	select {
	case <-ctx.Done():
		fmt.Println("Time out")
	case res := <-rep:
		fmt.Println(res)
		cancel()
	}
}
