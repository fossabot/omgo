#!/bin/bash

IPADDR=localhost
LOCALHOST=localhost
SERVICE_ROOT=backends
SERVICE_KIND=snowflake
SERVICE_NAME=snowflake-0
SERVICE_PORT=40001
MACHINE_ID=1
NETHOST=--net=host
ETCD_PORT=2379

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
    --port ${SERVICE_PORT} \
    --etcd-host http://${IPADDR}:${ETCD_PORT} \
    --service-root ${SERVICE_ROOT} \
    --service-kind ${SERVICE_KIND} \
    --service-name ${SERVICE_NAME} \
    --service-host ${IPADDR}
