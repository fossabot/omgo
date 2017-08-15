package main

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	endpoints = []string{
		"http://localhost:2379",
	}
)

func testSet() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer cli.Close()

	_, err = cli.Put(context.TODO(), "root/service_1/sub1", "s1_1")
	_, err = cli.Put(context.TODO(), "root/service_1/sub2", "s1_2")
	_, err = cli.Put(context.TODO(), "root/service_2/sub1", "s2_1")
	_, err = cli.Put(context.TODO(), "root/service_3/sub1", "s3_1")
	_, err = cli.Put(context.TODO(), "root/service_3/sub2", "s3_2")
	_, err = cli.Put(context.TODO(), "root/service_3/sub3", "s3_3")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := cli.Get(ctx, "root", clientv3.WithFromKey())
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}

func main() {
	names := []string{
		"service_1",
		"service_2",
		"service_3",
	}

	Init("root", endpoints, names)
	RegisterService("root/service_4/sub1", "s4_1")
}
