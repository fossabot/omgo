#!/bin/bash

LOCALHOST=127.0.0.1
IPADDR=127.0.0.1
SID=agent-0
NETHOST=--net=host

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

docker rm -f ${SID}
docker build --no-cache --rm=true -t agent .
docker run -d \
    --entrypoint /go/bin/agent \
    --name ${SID} \
    -e SERVICE_ID=${SID} \
    ${NETHOST} \
    -p 8888:8888 \
    -P agent \
    -l :8888 \
    -e http://${IPADDR}:2379 \
    -r backends \
    -s snowflake \
    -s game

# register service
curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/backends/agent/${SID} -d value=${IPADDR}:8888
