package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/utils"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	wg  sync.WaitGroup
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
