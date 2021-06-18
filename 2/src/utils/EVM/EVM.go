package EVM

import (
	"../Stack"
	"encoding/hex"
	"errors"
	"github.com/holiman/uint256"
	"golang.org/x/crypto/sha3"
	"strconv"
)

type EVM struct{
	Stack Stack.Stack
	Gas uint64
	Memory* Memory
}

func (e* EVM) DecodeInput(input string) error{

	length := len(input)
	if length == 0 {
		return nil
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

		w := e.Memory.Set(offset,value.Bytes())
		newCost := 3*w + w*w/512
		cost := newCost - e.Memory.lastGasCost
		e.Memory.lastGasCost = newCost
		e.Gas += cost


		e.DecodeInput(input[2:])

	} else if operation == "53" { //MSTORE8

		d,v := e.Stack.Pop(),e.Stack.Pop()
		offset,_ := strconv.ParseUint(d,16,64)
		value,_ := uint256.FromHex(Encode(v))

		w := e.Memory.Set8(offset,byte(value.Uint64()))
		newCost := 3*w + w*w/512
		cost := newCost - e.Memory.lastGasCost
		e.Memory.lastGasCost = newCost
		e.Gas += cost

		e.DecodeInput(input[2:])
	} else if operation == "01" { //ADD

		v1, v2 := e.Stack.Pop(), e.Stack.Pop()

		value1, _ := uint256.FromHex(Encode(v1))
		value2, _ := uint256.FromHex(Encode(v2))

		value1.Add(value1, value2)
		e.Stack.Push(value1.Hex()[2:])
		e.Gas += 3
		e.DecodeInput(input[2:])
	}else if operation == "02" { //MULL

		v1,v2  := e.Stack.Pop(),e.Stack.Pop()

		value1,_ := uint256.FromHex(Encode(v1))
		value2,_ := uint256.FromHex(Encode(v2))

		value1.Mul(value1,value2)
		e.Stack.Push(value1.Hex()[2:])
		e.Gas+=5
		e.DecodeInput(input[2:])
	}else if operation == "05" { //SDIV

		v1,v2  := e.Stack.Pop(),e.Stack.Pop()

		value1,_ := uint256.FromHex(Encode(v1))
		value2,_ := uint256.FromHex(Encode(v2))

		value1.SDiv(value1,value2)
		e.Stack.Push(value1.Hex()[2:])
		e.Gas+=5
		e.DecodeInput(input[2:])
	} else if operation == "0A" { //EXP

		b,ex  := e.Stack.Pop(),e.Stack.Pop()

		base,_ := uint256.FromHex(Encode(b))
		exp,_ := uint256.FromHex(Encode(ex))

		base.Exp(base,exp)
		e.Stack.Push(base.Hex()[2:])
		e.Gas+=50*uint64(len(exp.Bytes()))
		e.DecodeInput(input[2:])
	}else{
		return errors.New("wrong byte code")
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

	br :=0
	for string(str[br]) == "0"{
		br++
	}
	enc := "0x" + str[br:]
	return enc
}