package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/master-g/omgo/services"
	"github.com/master-g/omgo/utils"
)

func main() {
	defer utils.PrintPanicStack()
	test()
	services.Init("backends", []string{"http://localhost:2379"}, []string{"snowflake"})
}

type MyString struct {
	Length  int32
	Message [10]byte
}

type MyMessage struct {
	First   uint64
	Second  byte
	_       byte
	Third   uint32
	Message MyString
}

func test() {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &MyMessage{
		First:   100,
		Second:  0,
		Third:   100,
		Message: MyString{0, [10]byte{'H', 'e', 'l', 'l', 'o', '\n'}},
	})

	if err != nil {
		fmt.Printf("binary.Write failed:", err)
		return
	}

	msg := MyMessage{}
	err = binary.Read(buf, binary.LittleEndian, &msg)
	if err != nil {
		fmt.Printf("binary.Read failed:", err)
	}
	fmt.Println(msg)
}
