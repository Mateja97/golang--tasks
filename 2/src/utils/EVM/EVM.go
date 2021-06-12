package EVM

import (
	"../Stack"
	"github.com/holiman/uint256"
	"math"
	"strconv"
)

type EVM struct{
	Stack Stack.Stack
	Gas int
	Memory* Memory
}

func (e* EVM) DecodeInput(input string){

	length := len(input)
	if length == 0 {
		return
	}
	operation := input[0:2]

	if operation == "60" { //PUSH1
		e.Stack.Push(input[2:4])
		e.Gas += 3
		e.DecodeInput(input[4:])

	}else if operation == "61" { //PUSH2
		e.Stack.Push(input[2:6])
		e.Gas += 3
		e.DecodeInput(input[6:])

	}else if operation == "62" { //PUSH3
		e.Stack.Push(input[2:8])
		e.Gas += 3
		e.DecodeInput(input[8:])

	}else if operation == "7f" { //PUSH32
		e.Stack.Push(input[2:66])
		e.Gas += 3
		e.DecodeInput(input[66:])

	}else if operation == "52" { //MSTORE

		d,v := e.Stack.Pop(),e.Stack.Pop()
		offset,_ := strconv.ParseUint(d,16,64)
		value,_ := uint256.FromHex(Encode(v))

		e.Memory.Set(offset,value.Bytes())
		e.DecodeInput(input[2:])

	} else if operation == "53" { //MSTORE8

		d,v := e.Stack.Pop(),e.Stack.Pop()
		offset,_ := strconv.ParseUint(d,16,64)
		value,_ := uint256.FromHex(Encode(v))

		e.Memory.SetStore(offset,byte(value.Uint64()))
		e.DecodeInput(input[2:])
	} else if operation == "01" { //ADD

		v1, v2 := e.Stack.Pop(), e.Stack.Pop()

		value1, _ := uint256.FromHex(Encode(v1))
		value2, _ := uint256.FromHex(Encode(v2))
		//TO DO WITH SUM OVER 256 bits(big ints)
		value1.Add(value1, value2)
		e.Stack.Push(value1.Hex()[2:])
		e.Gas += 3
		e.DecodeInput(input[2:])
	}else if operation == "02" { //MULL

		v1,v2  := e.Stack.Pop(),e.Stack.Pop()

		value1,_ := uint256.FromHex(Encode(v1))
		value2,_ := uint256.FromHex(Encode(v2))
		//TO DO WITH PRODUCT OVER 256 bits(big ints)
		value1.Mul(value1,value2)
		e.Stack.Push(value1.Hex()[2:])
		e.Gas+=5
		e.DecodeInput(input[2:])
	}

}

func (e EVM) KECCAK256() string {
	return ""
}
//Adding 0x as prefix
func Encode(str string) string {

	br :=0
	for string(str[br]) == "0"{
		br++
	}
	enc := "0x" + str[br:]
	return enc
}
func WordSize(size uint64) uint64 {
	if size > math.MaxUint64-31 {
		return math.MaxUint64/32 + 1
	}
	return (size + 31) / 32
}