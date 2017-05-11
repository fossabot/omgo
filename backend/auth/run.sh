#!/bin/bash

IPADDR=127.0.0.1
LOCALHOST=127.0.0.1
SID=auth1
NETHOST=--net=host

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

# auth
docker rm -f ${SID}
docker build --no-cache --rm=true -t auth .
docker run --rm -d ${NETHOST} -p 40000:40000 \
    --name ${SID} \
    --entrypoint /go/bin/auth \
    auth \
    -l :40000 \
    -e http://${IPADDR}:2379

# register service
curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/backends/auth/${SID} -d value=${IPADDR}:40000
