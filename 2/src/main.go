package main

import (
	"../src/utils/EVM"
	"fmt"
)

func main() {

	var evm EVM.EVM// create a stack variable of type Stack
	evm = EVM.EVM{
		Memory: EVM.NewMemory(),
		Gas : 0 ,
	}
	evm.DecodeInput("6022600160200152")
	fmt.Println(evm.Memory.GetStore())

}