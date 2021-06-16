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
	evm.DecodeInput("60016020526002606452600361ff0052600362ffffff526005601053")
	//evm.DecodeInput("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00016000527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00026020527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff05604052")
	fmt.Println(evm.KECCAK256())
	fmt.Println(evm.Gas)
	//fmt.Println(evm.Memory.GetStore())
}