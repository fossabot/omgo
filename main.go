package main

import (
	"github.com/master-g/omgo/utils"
)

func main() {
	defer utils.PrintPanicStack()
	a := make([]byte, 1)
	a[3] = 123
}
