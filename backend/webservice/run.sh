#!/usr/bin/env bash

LOCALHOST=127.0.0.1
IPADDR=127.0.0.1
SERVICE_KIND=webservice
SERVICE_NAME=web-0
PORT=8080
NETHOST=--net=host

cp src/main/resources/config.json .

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     sed -i '' 's/localhost/'${IPADDR}'/g' config.json
     sed -i '' 's/v0.0.1/v0.0.1-autogen/g' config.json
     ;;
esac

# Database service
if [ "$1" = "rebuild" ]
then
    mvn -U clean package
    docker build --no-cache --rm=true -t ${SERVICE_KIND} .
fi

docker rm -f ${SERVICE_NAME}
docker run --rm -d ${NETHOST} -p ${PORT}:${PORT} \
    --name ${SERVICE_NAME} \
    ${SERVICE_KIND} \
