package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/kit/utils"
)

var (
	// increase while every agent goroutine was created
	// decrease while every agent goroutine was exited
	wg sync.WaitGroup
	// channel as signal to control agent goroutine to finish
	die = make(chan struct{})
)

// sigHandler handles unix signals
// this should be run in a goroutine
func sigHandler() {
	defer utils.PrintPanicStack()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)

	// loop until SIGTERM is received
	for {
		msg := <-ch
		switch msg {
		case syscall.SIGTERM:
			close(die) // notify all agent goroutine by close channel
			log.Info("sigterm received")
			log.Info("waiting for agents to close, please wait...")
			wg.Wait() // wait for all agent goroutine to finish
			log.Info("agents shutdown")
			os.Exit(0)
		}
	}
}
