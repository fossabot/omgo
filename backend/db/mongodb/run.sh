#!/bin/bash

# https://hub.docker.com/_/mongo/

SID=mongodb-0
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


# use master
# db.users.insert({usn:1})
# db.userExtra.insert({usn:1, secret:0})
# db.status.insert({usn:1, uid:1})
