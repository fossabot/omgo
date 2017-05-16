#!/bin/bash

# https://hub.docker.com/_/redis/

SID=redis-0
NETHOST=--net=host

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

docker rm -f ${SID}

docker run -d --name ${SID} ${NETHOST} \
  -p 6379:6379 \
  -v ${PWD}/redis.conf:/usr/local/etc/redis/redis.conf \
  -v ${PWD}/db:/data/ \
  redis redis-server \
  /usr/local/etc/redis/redis.conf
