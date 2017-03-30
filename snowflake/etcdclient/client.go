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

func KeysAPI() etcdclient.KeysAPI {
	return etcdclient.NewKeysAPI(client)
}
