package EVM

import (
	"encoding/hex"
	"errors"

	"github.com/holiman/uint256"
	"golang.org/x/2/src/utils/Stack"
	"golang.org/x/crypto/sha3"
)

type EVM struct {
	Stack  Stack.Stack
	Gas    uint64
	Memory *Memory
}

var evm *EVM

//Make EVM Class as singleton
func init() {
	evm = &EVM{
		Memory: NewMemory(),
		Gas:    0,
	}
}
func GetInstance() *EVM {
	return evm
}

var Actions = map[string]interface{}{
	"60": Push1,
	"61": Push2,
	"62": Push3,
	"7f": Push32,
	"52": Mstore,
	"53": Mstore8,
	"01": Add,
	"02": Mul,
	"05": SDiv,
	"0a": Exp,
}

func (e *EVM) DecodeInput(input string) error {

	length := len(input)
	if length == 0 {
		return nil
	}
	operation := input[0:2]
	action, ok := Actions[operation]

	if ok {
		offset := action.(func(*EVM, string) int)(e, input)
		_ = e.DecodeInput(input[offset:])

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

	i := 0
	for i < len(str) && string(str[i]) == "0" {
		i++
	}
	enc := "0x"
	if i == len(str) {
		enc += "0"

	} else {
		enc += str[i:]
	}
	return enc
}
func Push1(e *EVM, input string) int {

	e.Stack.Push(input[2:4])
	e.Gas += 3
	return 4
}
func Push2(e *EVM, input string) int {
	e.Stack.Push(input[2:6])
	e.Gas += 3
	return 6
}
func Push3(e *EVM, input string) int {
	e.Stack.Push(input[2:8])
	e.Gas += 3
	return 8
}
func Push32(e *EVM, input string) int {
	e.Stack.Push(input[2:66])
	e.Gas += 3
	return 66
}
func Mstore(e *EVM, input string) int {
	d, v := e.Stack.Pop(), e.Stack.Pop()
	offset, _ := uint256.FromHex(Encode(d))
	value, _ := uint256.FromHex(Encode(v))

	//calculate gas by formula
	w := e.Memory.Set(offset.Uint64(), value.Bytes())
	newCost := 3*w + w*w/512
	cost := newCost - e.Memory.lastGasCost
	e.Memory.lastGasCost = newCost
	e.Gas += 3 + cost
	return 2

}

func Mstore8(e *EVM, input string) int {
	d, v := e.Stack.Pop(), e.Stack.Pop()
	offset, _ := uint256.FromHex(Encode(d))
	value, _ := uint256.FromHex(Encode(v))

	//calculate gas by formula
	w := e.Memory.Set8(offset.Uint64(), byte(value.Uint64()))
	newCost := 3*w + w*w/512
	cost := newCost - e.Memory.lastGasCost
	e.Memory.lastGasCost = newCost
	e.Gas += 3 + cost
	return 2
}
func Add(e *EVM, input string) int {
	v1, v2 := e.Stack.Pop(), e.Stack.Pop()

	value1, _ := uint256.FromHex(Encode(v1))
	value2, _ := uint256.FromHex(Encode(v2))

	value1.Add(value1, value2)
	e.Stack.Push(value1.Hex()[2:])
	e.Gas += 3
	return 2
}
func Mul(e *EVM, input string) int {
	v1, v2 := e.Stack.Pop(), e.Stack.Pop()

	value1, _ := uint256.FromHex(Encode(v1))
	value2, _ := uint256.FromHex(Encode(v2))

	value1.Mul(value1, value2)
	e.Stack.Push(value1.Hex()[2:])
	e.Gas += 5
	return 2
}
func SDiv(e *EVM, input string) int {
	v1, v2 := e.Stack.Pop(), e.Stack.Pop()

	value1, _ := uint256.FromHex(Encode(v1))
	value2, _ := uint256.FromHex(Encode(v2))

	value1.SDiv(value1, value2)
	e.Stack.Push(value1.Hex()[2:])
	e.Gas += 5
	return 2
}

func Exp(e *EVM, input string) int {
	b, ex := e.Stack.Pop(), e.Stack.Pop()

	base, _ := uint256.FromHex(Encode(b))
	exp, _ := uint256.FromHex(Encode(ex))

	base.Exp(base, exp)
	e.Stack.Push(base.Hex()[2:])
	e.Gas += 10 + 50*uint64(len(exp.Bytes()))
	return 2
}
