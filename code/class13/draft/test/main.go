package main

import "fmt"

func main() {
	var a = 2
	var b = 3
	res := calculate(a, b)
	fmt.Println(res)
	return
}

func calculate(x, y int) int {
	return x * y
}
