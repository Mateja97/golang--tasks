package EVM

import (
	"encoding/hex"
	"errors"
	"strconv"

	"../Stack"
	"github.com/holiman/uint256"
	"golang.org/x/crypto/sha3"
)

type EVM struct {
	Stack  Stack.Stack
	Gas    uint64
	Memory *Memory
}

func (e *EVM) DecodeInput(input string) error {

	length := len(input)
	if length == 0 {
		return nil
	}
	operation := input[0:2]

	if operation == "60" { //PUSH1
		e.Push1(input)

	} else if operation == "61" { //PUSH2
		e.Push2(input)

	} else if operation == "62" { //PUSH3
		e.Push3(input)

	} else if operation == "7f" { //PUSH32
		e.Push32(input)

	} else if operation == "52" { //MSTORE
		e.Mstore()
		e.DecodeInput(input[2:])

	} else if operation == "53" { //MSTORE8

		e.Mstore8()
		e.DecodeInput(input[2:])

	} else if operation == "01" { //ADD

		e.Add()
		e.DecodeInput(input[2:])

	} else if operation == "02" { //MULL

		e.Mul()
		e.DecodeInput(input[2:])

	} else if operation == "05" { //SDIV

		e.SDiv()
		e.DecodeInput(input[2:])

	} else if operation == "0A" { //EXP

		e.Exp()
		e.DecodeInput(input[2:])
	} else {
		return errors.New("Wrong byte code")
	}
	return nil

}

func (e EVM) KECCAK256() string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(e.Memory.store)
	buf := hash.Sum(nil)
	return hex.EncodeToString(buf)
}

//Adding 0x as prefix
func Encode(str string) string {

	br := 0
	for string(str[br]) == "0" {
		br++
	}
	enc := "0x" + str[br:]
	return enc
}
func (e *EVM) Push1(input string) {

	e.Stack.Push(input[2:4])
	e.Gas += 3
	e.DecodeInput(input[4:])
}
func (e *EVM) Push2(input string) {
	e.Stack.Push(input[2:6])
	e.Gas += 3
	e.DecodeInput(input[6:])
}
func (e *EVM) Push3(input string) {
	e.Stack.Push(input[2:8])
	e.Gas += 3
	e.DecodeInput(input[8:])
}
func (e *EVM) Push32(input string) {
	e.Stack.Push(input[2:66])
	e.Gas += 3
	e.DecodeInput(input[66:])
}
func (e *EVM) Mstore() {
	d, v := e.Stack.Pop(), e.Stack.Pop()
	offset, _ := strconv.ParseUint(d, 16, 64)
	value, _ := uint256.FromHex(Encode(v))

	//calculate gas by formula
	w := e.Memory.Set(offset, value.Bytes())
	newCost := 3*w + w*w/512
	cost := newCost - e.Memory.lastGasCost
	e.Memory.lastGasCost = newCost
	e.Gas += cost

}

func (e *EVM) Mstore8() {
	d, v := e.Stack.Pop(), e.Stack.Pop()
	offset, _ := strconv.ParseUint(d, 16, 64)
	value, _ := uint256.FromHex(Encode(v))
	//calculate gas by formula
	w := e.Memory.Set8(offset, byte(value.Uint64()))
	newCost := 3*w + w*w/512
	cost := newCost - e.Memory.lastGasCost
	e.Memory.lastGasCost = newCost
	e.Gas += cost
}
func (e *EVM) Add() {
	v1, v2 := e.Stack.Pop(), e.Stack.Pop()

	value1, _ := uint256.FromHex(Encode(v1))
	value2, _ := uint256.FromHex(Encode(v2))

	value1.Add(value1, value2)
	e.Stack.Push(value1.Hex()[2:])
	e.Gas += 3
}
func (e *EVM) Mul() {
	v1, v2 := e.Stack.Pop(), e.Stack.Pop()

	value1, _ := uint256.FromHex(Encode(v1))
	value2, _ := uint256.FromHex(Encode(v2))

	value1.Mul(value1, value2)
	e.Stack.Push(value1.Hex()[2:])
	e.Gas += 5
}
func (e *EVM) SDiv() {
	v1, v2 := e.Stack.Pop(), e.Stack.Pop()

	value1, _ := uint256.FromHex(Encode(v1))
	value2, _ := uint256.FromHex(Encode(v2))

	value1.SDiv(value1, value2)
	e.Stack.Push(value1.Hex()[2:])
	e.Gas += 5
}

func (e *EVM) Exp() {
	b, ex := e.Stack.Pop(), e.Stack.Pop()

	base, _ := uint256.FromHex(Encode(b))
	exp, _ := uint256.FromHex(Encode(ex))

	base.Exp(base, exp)
	e.Stack.Push(base.Hex()[2:])
	e.Gas += 50 * uint64(len(exp.Bytes()))
}
