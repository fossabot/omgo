/*
MIT License

Copyright (c) 2017 Master.G

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package etcdclient

import (
	log "github.com/Sirupsen/logrus"
	etcdclient "github.com/coreos/etcd/client"
)

var client etcdclient.Client

func Init(endpoints []string) {
	// config
	cfg := etcdclient.Config{
		Endpoints: endpoints,
		Transport: etcdclient.DefaultTransport,
	}

	// create client
	c, err := etcdclient.New(cfg)
	if err != nil {
		log.Error(err)
		return
	}
	client = c
}

// KeysAPI builds a KeysAPI that interacts with etcd's key-value
// API over HTTP
func KeysAPI() etcdclient.KeysAPI {
	return etcdclient.NewKeysAPI(client)
}

// NewOptions builds an empty etcd GetOptions
func NewOptions() etcdclient.GetOptions {
	return etcdclient.GetOptions{}
}

// NewWatcherOptions builds a new etcd WatcherOptions
func NewWatcherOptions(recursive bool) *etcdclient.WatcherOptions {
	return &etcdclient.WatcherOptions{Recursive: recursive}
}
