package main

import (
	"fmt"
	"github.com/master-g/omgo/mt19937"
	"github.com/master-g/omgo/utils"
	"time"
)

func main() {
	defer utils.PrintPanicStack(time.Now())
	ctx := mt19937.NewContext(0)
	fmt.Println(ctx.NextInt32())
	a := make([]int, 3)
	a[3] = 4
}
