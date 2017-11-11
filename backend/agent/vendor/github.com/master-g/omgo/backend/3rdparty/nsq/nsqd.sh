#!/usr/bin/env bash

IPADDR=localhost
LOCALHOST=localhost
NSQD_SID=nsqd
NETHOST=--net=host

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

#nsqlookupd address via etcd
LOOKUPD_HOST=$(ETCDCTL_API=3 etcdctl get backends/nsq/nsqlookupd | sed -n 2p)

#nsqd

if [ "$1" = "rebuild" ]
then
    docker pull nsqio/nsq
fi

docker rm -f ${NSQD_SID}
docker run --rm -d ${NETHOST} -p 4150:4150 -p 4151:4151 \
    --name ${NSQD_SID} \
    nsqio/nsq /nsqd \
    --broadcast-address=${LOOKUPD_HOST} \
    --lookupd-tcp-address=${LOOKUPD_HOST}:4160
