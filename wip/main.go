package main

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/prometheus/common/log"
	"golang.org/x/net/context"
	"time"
)

func main() {
	endpoints := []string{
		"http://localhost:2379",
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer cli.Close()

	_, err = cli.Put(context.TODO(), "holyshit/sub1", "ass")
	_, err = cli.Put(context.TODO(), "holyshit/sub2", "fuck")
	_, err = cli.Put(context.TODO(), "holyshit/sub3", "shit")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := cli.Get(ctx, "holyshit", clientv3.WithFromKey())
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}
