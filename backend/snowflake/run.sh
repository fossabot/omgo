#!/bin/bash

IPADDR=localhost
LOCALHOST=localhost
SERVICE_NAME=snowflake
SERVICE_PORT=40001
MACHINE_ID=1
SNOWFLAKE_SID=snowflake-0
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
    docker build --no-cache --rm=true -t ${SERVICE_NAME} .
fi

docker rm -f ${SNOWFLAKE_SID}
docker run --rm -d ${NETHOST} -p ${SERVICE_PORT}:${SERVICE_PORT} \
    --name ${SNOWFLAKE_SID} \
    -e SERVICE_ID=${SNOWFLAKE_SID} \
    -e MACHINE_ID=${MACHINE_ID} \
    --entrypoint /go/bin/${SERVICE_NAME} \
    ${SERVICE_NAME} \
    --service-key backends/${SERVICE_NAME}/${SNOWFLAKE_SID} \
    --service-host ${IPADDR} \
    -p ${SERVICE_PORT} \
    -e http://${IPADDR}:2379
