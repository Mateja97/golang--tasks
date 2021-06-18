package main

import (
	"../src/utils/EVM"
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {

	file, err := os.Open("data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var evm EVM.EVM// create a stack variable of type Stack
		evm = EVM.EVM{
			Memory: EVM.NewMemory(),
			Gas : 0 ,
		}
		input := scanner.Text()
		fmt.Println("input: ",input)
		if len(input) %2 != 0 {
			fmt.Println("Wrong input, code should be dividable by 2")
			continue
		}
		err:= evm.DecodeInput(input)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Keccak: ",evm.KECCAK256())
		fmt.Println("Gas consumed", evm.Gas)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}