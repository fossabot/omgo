#!/bin/bash

# https://hub.docker.com/_/mongo/

SID=mongodb1
NETHOST=--net=host

case "$(uname -s)" in
   Darwin)
     IPADDR=$(ifconfig en0 | grep "inet " | cut -d " " -f2)
     NETHOST=''
     ;;
esac

docker rm -f ${SID}

docker run -d --name ${SID} ${NETHOST} \
  -p 37017:27017 \
  -v ${PWD}/mongod.conf:/etc/mongod.conf \
  -v ${PWD}/db:/data/db \
  mongo \
  --config /etc/mongod.conf
