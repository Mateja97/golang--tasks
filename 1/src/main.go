package main

import (
    "bufio"
    "encoding/hex"
    "fmt"
    "log"
    "os"
    "strconv"
)
func rlpDecode(input string,output  *[]string) {
    if len(input) == 0{
        return //out of recursion
    }
    dtype,offSet,dataLength:= decodeLength(input) //calculating starting offset of data and datalength, dtype is checking is it string,list or errr
    if dtype == "str"{
        *output = append(*output,input[offSet:offSet+dataLength]) // add found string to output result
        rlpDecode(input[offSet+dataLength:],output) // recursive call without found string

    }else if dtype == "list"{
        *output = append(*output,"5b") //adding hex value for start bracket
        rlpDecode(input[offSet:],output) // recursive call to iterate through list
        *output = append(*output,"5d") //adding hex value for end bracket
    }else if dtype == "error"{
        *output = append(make([]string,1),hex.EncodeToString([]byte("wrong rlp coding")))
    }
}
func decodeLength(input string) (string,int64,int64){
    length := int64(len(input))
    if length == 0{
        return "error",1,1
    }
    //Validate input
    if !isHex(input[0:2]) || length %2 != 0 {
        return "error",1,1
    }
    prefix, _ := strconv.ParseInt(input[0:2], 16, 64)// taking first byte to check what type of data is in input
    
    if prefix <= 0x7f {//input is byte itself
        return "str",0,1
    }else if prefix <= 0xb7 && length >= 2*(prefix - 0x80){//short string

        strLen := prefix - 0x80
        return "str",2,strLen*2

    } else if prefix <= 0xbf && length >= 2*(prefix-0xb7) { //long string

        offset := 2*(prefix - 0xb7)
        l, _ := strconv.ParseInt(input[2:2*offset], 16, 64) // taking length of bytes for string bigger then 55 bytes
        if length > (prefix - 0xb7 + l)*2{
            return "str", 2*offset, l*2
        }
        return "error",1,1
    }else if prefix <= 0xf7 { //short list
        listLen:= prefix - 0xc0
        return "list",2,listLen*2
    }else if prefix <= 0xff { // long list
        offset := 2*(prefix - 0xf7)
        l, _ := strconv.ParseInt(input[2:2*offset], 16, 64) 
        if length > (prefix - 0xf7 + l)*2{
            return "list", 2+2*offset, l*2
        }
        return "error",1,1

    } else{
        return "error",1,1
    }

}
//Function to verify is input hex encoding
func isHex(content string) bool {
    isHex := true
    _, err := strconv.ParseUint(content, 16, 64)
    if err != nil {
        isHex = false
    }
    return isHex
}
func main() {

    file, err := os.Open("../data.txt")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        var decodings []string
        rlpDecode(scanner.Text(),&decodings)
        fmt.Print("Result: ")
        for _,elem := range decodings{
            out,_ :=hex.DecodeString(elem)
            fmt.Print(string(out)+ " ")
        }
        fmt.Print("\n")
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

}