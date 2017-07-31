#!/usr/bin/env bash

LOCALHOST=127.0.0.1
IPADDR=127.0.0.1
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
mvn clean package

docker rm -f ${SID}
docker build --no-cache --rm=true -t ${SERVICE_NAME} .
docker run --rm -d ${NETHOST} -p ${PORT}:${PORT} \
    --name ${SID} \
    ${SERVICE_NAME} \

# register service
curl -q -L -X PUT http://${LOCALHOST}:2379/v2/keys/backends/${SERVICE_NAME}/${SID} -d value=${IPADDR}:${PORT}
