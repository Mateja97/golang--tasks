package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"golang.org/x/2/src/utils/EVM"
)

//Load input as string array
func inputData() []string {
	file, err := os.Open("../data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var data []string
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return data
}
func main() {

	data := inputData()
	//iterate through input
	for _, input := range data {
		evm := EVM.GetInstance()

		fmt.Println("input:", input)

		if len(input)%2 != 0 {
			fmt.Println("Wrong input, code should be dividable by 2")
			fmt.Println("")
			continue
		}
		err := evm.DecodeInput(input)
		if err != nil {
			fmt.Println(err)
			fmt.Println("")
			continue
		}
		fmt.Println("Keccak:", evm.KECCAK256())
		fmt.Println("Gas consumed:", evm.Gas)
		fmt.Println("")
		//Reset EVM for new input
		evm.Gas = 0
		evm.Memory = EVM.NewMemory()
	}

}
