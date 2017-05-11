#!/bin/bash

IPADDR=127.0.0.1
LOCALHOST=127.0.0.1
SID=reception-0
NETHOST=--net=host

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

# auth
docker rm -f ${SID}
docker build --no-cache --rm=true -t reception .
docker run --rm -d ${NETHOST} \
    -p 8080:8080 \
    --name ${SID} \
    --entrypoint /go/bin/reception \
    reception \
    -l :8080 \
    -e http://${IPADDR}:2379
