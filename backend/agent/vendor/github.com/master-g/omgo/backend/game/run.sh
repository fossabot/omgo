#!/bin/bash

LOCALHOST=127.0.0.1
IPADDR=127.0.0.1
SID=game-0
NETHOST=--net=host

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

docker rm -f ${SID}
docker build --no-cache --rm=true -t game .
docker run -d \
    --entrypoint /go/bin/game \
    --name ${SID} \
    -e SERVICE_ID=${SID} \
    ${NETHOST} \
    -p 10000:10000 \
    -P game \
    -l :10000 \
    -i game-0
    -e http://${IPADDR}:2379 \
    -r backends \
    -s snowflake \
    -s dbservice

# register service
curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/backends/game/${SID} -d value=${IPADDR}:10000