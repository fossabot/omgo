package main

import (
	"fmt"
	"github.com/master-g/omgo/services"
	"github.com/master-g/omgo/utils"
)

func main() {
	defer utils.PrintPanicStack()
	test()
	services.Init("backends", []string{"http://localhost:2379"}, []string{"snowflake"})
}

func test() {
	a := make([]byte, 12)
	b := []byte{0, 1, 2, 3}
	a = append(a, b...)
	fmt.Println(a)
}
