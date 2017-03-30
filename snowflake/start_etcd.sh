#!/usr/bin/env bash

# --net=host, bind the container to the host interface
# so other container can access this container via localhost:port
# ps. --net=host is not working on macOS, see https://github.com/docker/for-mac/issues/68

IPADDR=127.0.0.1

docker rm -f etcd
docker run -d --net=host -p 2379:2379 -p 2380:2380 --name etcd quay.io/coreos/etcd \
    /usr/local/bin/etcd \
    --data-dir=data.etcd --name etcd0 \
    --initial-advertise-peer-urls http://$IPADDR:2380 --listen-peer-urls http://$IPADDR:2380 \
    --advertise-client-urls http://0.0.0.0:2379 --listen-client-urls http://0.0.0.0:2379 \
    --initial-cluster etcd0=http://$IPADDR:2380 \
    --initial-cluster-state new --initial-cluster-token my-etcd-token
