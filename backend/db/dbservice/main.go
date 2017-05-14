package main

import (
    log "github.com/Sirupsen/logrus"
    "github.com/master-g/omgo/utils"
)

const (
    defaultETCD = "http://127.0.0.1:2379"
)

func main() {
    log.SetLevel(log.DebugLevel)
    defer utils.PrintPanicStack()


}