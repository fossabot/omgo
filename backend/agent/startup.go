package main

import "github.com/master-g/omgo/services"

func startup(root string, hosts, names []string) {
	go sigHandler()
	services.Init(root, hosts, names)
}
