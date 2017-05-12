#!/bin/bash

IPADDR=127.0.0.1
LOCALHOST=127.0.0.1
SERVICE_NAME=dbservice
SID=dbs-0
PORT=60001
NETHOST=--net=host

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

# Database service
docker rm -f ${SID}
docker build --no-cache --rm=true -t ${SERVICE_NAME} .
docker run --rm -d ${NETHOST} -p ${PORT}:${PORT} \
    --name ${SID} \
    -e SERVICE_ID=${SID} \
    -e MACHINE_ID=1 \
    --entrypoint /go/bin/${SERVICE_NAME} \
    ${SERVICE_NAME} \
    -p ${PORT} \
    -e http://${IPADDR}:2379

# register service
curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/backends/${SERVICE_NAME}/${SID} -d value=${IPADDR}:${PORT}
