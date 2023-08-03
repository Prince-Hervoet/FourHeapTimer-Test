package main

import (
	"fmt"
	"fourheap/core"
	"time"
)

func main() {
	test := core.NewTimer(func(a any) any {
		fmt.Println("I am ran1.hahaha")
		return nil
	}, nil, 5000)

	test2 := core.NewTimer(func(a any) any {
		fmt.Println("I am ran2.hahaha")
		return nil
	}, nil, 1000)

	test3 := core.NewTimer(func(a any) any {
		fmt.Println("I am ran3.hahah")
		return nil
	}, nil, 3000)

	fh := core.NewFourHeap(1000)
	fh.AddTimer(test)
	fh.AddTimer(test2)
	fh.AddTimer(test3)
	fh.Start()

	for {
		time.Sleep(1 * time.Second)
		fmt.Println(fh.GetSize())
	}
}
