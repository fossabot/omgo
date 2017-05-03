#!/bin/bash
IPADDR=127.0.0.1
NETHOST=--net=host

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

export DOCKER_HOST_IP=${IPADDR}

echo '==> building environment'

docker-compose build --pull

echo '==> launching environment'

docker-compose up -d
