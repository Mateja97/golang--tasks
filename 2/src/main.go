package main

import (
	"../src/utils/Stack"
	"fmt"
)

func main() {

	var stack Stack.Stack// create a stack variable of type Stack

	stack.Push("this")
	stack.Push("is")
	stack.Push("sparta!!")

	for len(stack.Values) > 0 {
		x, y := stack.Pop()
		if y == true {
			fmt.Println(x)
		}
	}
}