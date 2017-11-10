#!/usr/bin/env bash

IPADDR=localhost
LOCALHOST=localhost
NSQLOOKUPD_SID=lookupd
NETHOST=--net=host

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

#nsqlookupd

if [ "$1" = "rebuild" ]
then
    docker pull nsqio/nsq
fi

docker rm -f ${NSQLOOKUPD_SID}
docker run --rm -d ${NETHOST} -p 4160:4160 -p 4161:4161 \
    --name ${NSQLOOKUPD_SID} \
    nsqio/nsq /nsqlookupd

# register to ETCD
ETCDCTL_API=3 etcdctl put backends/nsq/nsqlookupd ${IPADDR}