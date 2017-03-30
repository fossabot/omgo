package main

import (
	"fmt"
	"github.com/master-g/omgo/mt19937"
)

func main() {
	ctx := mt19937.NewContext(0)
	fmt.Println(ctx.NextInt32())
}
