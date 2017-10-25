#!/usr/bin/env bash

IPADDR=localhost
LOCALHOST=localhost
ETCD_SID=etcd0
NETHOST=--net=host

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

#Get ETCD IP address
#ETCD_IP=$(docker inspect etcd | python -c 'import json,sys;obj=json.load(sys.stdin);print obj[0]["NetworkSettings"]["IPAddress"]')

#ETCD

if [ "$1" = "rebuild" ]
then
    docker pull quay.io/coreos/etcd
fi

docker rm -f etcd
docker run --rm -d ${NETHOST} -p 2379:2379 -p 2380:2380 \
    --name etcd quay.io/coreos/etcd \
    /usr/local/bin/etcd \
    --name ${ETCD_SID} \
    --data-dir=data.etcd \
    --initial-advertise-peer-urls http://${LOCALHOST}:2380 \
    --listen-peer-urls http://${LOCALHOST}:2380 \
    --advertise-client-urls http://0.0.0.0:2379 \
    --listen-client-urls http://0.0.0.0:2379 \
    --initial-cluster ${ETCD_SID}=http://${LOCALHOST}:2380 \
    --initial-cluster-state new \
    --initial-cluster-token my-etcd-token