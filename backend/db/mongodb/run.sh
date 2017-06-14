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

# *** CONNECT TO MONGODB ***
# docker exec -it mongodb-0 mongo

# *** CREATE ADMIN DB AND ADMIN USER ***
# use admin
# db.createUser({ user:'admin', pwd:'admin', roles:['userAdminAnyDatabase', 'dbAdminAnyDatabase', 'readWriteAnyDatabase']})

# *** CREATE DB AND USER FOR DB CLIENT ***
# use master
# db.createUser({user:'dbclient', pwd:'12345678', roles:['dbOwner']})

# *** INITIALIZE DB ***
# db.status.insert({key:'user', usn:10000, uid:10000})

# docker exec -it mongodb-0 mongo admin -u admin -p admin
# docker exec -it mongodb-0 mongo master -u dbclient -p 12345678