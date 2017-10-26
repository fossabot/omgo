#!/bin/bash

IPADDR=localhost
LOCALHOST=localhost
SERVICE_KIND=snowflake
SERVICE_PORT=40001
MACHINE_ID=1
SERVICE_NAME=snowflake-0
NETHOST=--net=host

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

#Snowflake
if [ "$1" = "rebuild" ]
then
    docker build --no-cache --rm=true -t ${SERVICE_KIND} .
fi

docker rm -f ${SERVICE_NAME}
docker run --rm -d ${NETHOST} -p ${SERVICE_PORT}:${SERVICE_PORT} \
    --name ${SERVICE_NAME} \
    -e SERVICE_ID=${SERVICE_NAME} \
    -e MACHINE_ID=${MACHINE_ID} \
    --entrypoint /go/bin/${SERVICE_KIND} \
    ${SERVICE_KIND} \
    --service-key backends/${SERVICE_KIND}/${SERVICE_NAME} \
    --service-host ${IPADDR} \
    --port ${SERVICE_PORT} \
    --etcd http://${IPADDR}:2379
