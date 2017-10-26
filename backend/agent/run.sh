#!/bin/bash

# service configuration
IPADDR=127.0.0.1
SERVICE_ROOT=backends
SERVICE_KIND=agent
SERVICE_NAME=agent-0
SERVICE_PORT=30001
AGENT_PORT=8888
SERVICE_DATASERVICE=dataservice
SERVICE_GAMESERVICE=game
ETCD_PORT=2379

# --net=host does not work properly on OSX
# since docker runs in a virtual machine on OSX
# omit --net=host and use OSX's real IP address can bypass this problem

NETHOST=--net=host
case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

if [ "$1" = "rebuild" ]
then
    docker build --no-cache --rm=true -t ${SERVICE_KIND} .
fi

docker rm -f ${SERVICE_NAME}
docker run --rm -d ${NETHOST} -p ${AGENT_PORT}:${AGENT_PORT} \
    --name ${SERVICE_NAME} \
    -e SERVICE_ID=${SERVICE_NAME} \
    --entrypoint /go/bin/${SERVICE_KIND} \
    ${SERVICE_KIND} \
    -P ${SERVICE_KIND} \
    -l ${IPADDR}:${AGENT_PORT} \
    -k ${SERVICE_KIND} \
    -n ${SERVICE_NAME} \
    -e http://${IPADDR}:${ETCD_PORT} \
    -r ${SERVICE_ROOT} \
    -s ${SERVICE_DATASERVICE} \
    -s ${SERVICE_GAMESERVICE}
