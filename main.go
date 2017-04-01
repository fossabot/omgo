package main

import (
	"github.com/master-g/omgo/services"
	"github.com/master-g/omgo/utils"
)

func main() {
	defer utils.PrintPanicStack()
	services.Init("backends", []string{"http://localhost:2379"}, []string{"snowflake"})
}
