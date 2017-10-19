package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/utils"
)

var (
	// increase while every agent goroutine was created
	// decrease while every agent goroutine was exited
	wg sync.WaitGroup
	// channel as signal to control agent goroutine to finish
	die = make(chan struct{})
)

// handle unix signals
func sigHandler() {
	defer utils.PrintPanicStack()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)

	for {
		msg := <-ch
		switch msg {
		case syscall.SIGTERM:
			close(die)
			log.Info("sigterm received")
			log.Info("waiting for agents to close, please wait...")
			wg.Wait()
			log.Info("agent shutdown")
			os.Exit(0)
		}
	}
}
