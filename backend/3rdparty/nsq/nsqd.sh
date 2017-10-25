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

#nsqd

if [ "$1" = "rebuild" ]
then
    docker pull nsqio/nsq
fi

docker rm -f ${NSQD_SID}
docker run --rm -d ${NETHOST} -p 4150:4150 -p 4151:4151 \
    --name ${NSQD_SID} \
    nsqio/nsq /nsqd \
    --broadcast-address=<host> \
    --lookupd-tcp-address=<host>:<port>
    